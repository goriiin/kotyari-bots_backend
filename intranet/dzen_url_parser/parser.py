from __future__ import annotations

import logging
import time
from typing import Set, List, Optional, Dict, Iterable
from urllib.parse import urlparse, urlunparse

from bs4 import BeautifulSoup
from selenium.webdriver.remote.webdriver import WebDriver

from .config import settings
from .redis_adapter import LinkStorer, LinkItem

logger = logging.getLogger(__name__)


def normalize_dzen_article_url(href: str) -> Optional[str]:
    """
    Normalize Dzen article URLs to the canonical https://dzen.ru/a/<id> form.
    """
    try:
        p = urlparse(href)
        if p.scheme in ("http", "https") and p.netloc.endswith("dzen.ru"):
            if p.path.startswith("/a/"):
                clean = p._replace(query="", fragment="")
                return urlunparse(clean.geturl())
        if href.startswith("/a/"):
            return f"https://dzen.ru{href.split('?')[0]}"
        return None
    except Exception as e:
        logger.error("normalize error: %s", e)
        return None


def collect_links_from_html(html: str) -> List[str]:
    """
    Extract article links from Dzen listing page HTML.
    """
    soup = BeautifulSoup(html, "html.parser")
    found: List[str] = []

    # Primary selectors
    a_candidates = []
    for dtid in ("card-article-link", "card-article-title-link"):
        a_candidates.extend(soup.find_all("a", attrs={"data-testid": dtid}))

    # Additional patterns observed in feed
    for article in soup.find_all(attrs={"data-testid": "floor-image-card"}):
        a = article.find("a", attrs={"data-testid": "card-article-link"})
        if a is not None:
            a_candidates.append(a)

    for wrap in soup.find_all(attrs={"data-testid": "card-part-title"}):
        parent = wrap.parent
        if parent:
            a = parent.find("a", attrs={"data-testid": "card-article-title-link"})
            if a is not None:
                a_candidates.append(a)

    for a in a_candidates:
        href = a.get("href")
        if not href:
            continue
        norm = normalize_dzen_article_url(href)
        if norm:
            found.append(norm)

    # Fallback: scan all <a href=...>
    if not found:
        for a in soup.find_all("a", href=True):
            norm = normalize_dzen_article_url(a["href"])
            if norm:
                found.append(norm)

    # unique, keep order
    seen: Set[str] = set()
    uniq: List[str] = []
    for url in found:
        if url not in seen:
            uniq.append(url)
            seen.add(url)

    if settings.DEBUG_LOG_SELECTORS:
        preview = ", ".join(uniq[:3])
        logger.info("parser.debug collected=%d preview=%s", len(uniq), preview)

    return uniq


def scroll_page(driver: WebDriver, count: int, delay_seconds: int) -> None:
    for i in range(count):
        driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
        time.sleep(delay_seconds)
        if settings.DEBUG_LOG_SELECTORS:
            print("parser.debug Scrolling down... %d/%d", i + 1, count)


def extract_category_from_topic_url(u: str) -> Optional[str]:
    p = urlparse(u)
    path = (p.path or "/").rstrip("/")
    if not path or path == "/":
        return None
    # /topic/<category>
    parts = path.split("/")
    return parts[-1] if parts else None


def publish_batch(linkstorer: Optional[LinkStorer], urls: Iterable[str], category: Optional[str]) -> int:
    """
    Build [{"url","category"}] and delegate to LinkStorer.
    """
    if not linkstorer:
        return 0
    items: List[LinkItem] = [{"url": u, "category": category} for u in urls]
    return linkstorer.store_links(items)


def parse_dzen_for_links_with_category(driver: WebDriver, link_storer: Optional[LinkStorer] = None) -> List[Dict[str, Optional[str]]]:
    """
    Iterate configured topic URLs, scroll, collect links and publish them with category.
    Returns the deduplicated mapping converted to list of {"url","category"}.
    """
    print("Starting Dzen parsing process for all configured URLs...")

    unique: Dict[str, Optional[str]] = {}
    published_seen: Set[str] = set()

    try:
        for topic_url in settings.DZEN_ARTICLES_URLs:
            category = extract_category_from_topic_url(topic_url)

            # Navigate and scroll the feed
            driver.get(topic_url)
            time.sleep(5)
            for _ in range(settings.DZEN_SCROLL_COUNT):
                driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
                time.sleep(settings.DZEN_SCROLL_DELAY_SECONDS)

            # Collect links from HTML
            urls = collect_links_from_html(driver.page_source)

            # Merge into unique and prepare batch to publish
            new_urls = [u for u in urls if u not in unique]
            for u in new_urls:
                unique[u] = category

            to_publish = [u for u in new_urls if u not in published_seen]
            if to_publish and link_storer:
                published = publish_batch(link_storer, to_publish, category)
                if published > 0:
                    published_seen.update(to_publish)

        return [{"url": u, "category": c} for u, c in unique.items()]
    except Exception as e:
        logger.error("parser error during parsing: %s", e)
        return [{"url": u, "category": c} for u, c in unique.items()]
