from __future__ import annotations

import time
from typing import Set, List, Optional, Dict, Iterable
from urllib.parse import urlparse

from bs4 import BeautifulSoup
from selenium.webdriver.remote.webdriver import WebDriver

from .config import settings
from .redis_adapter import LinkStorer, LinkItem

def normalize_dzen_article_url(href: str) -> Optional[str]:
    """
    Normalize Dzen article URLs to the canonical https://dzen.ru/a/<id> form.
    """
    try:
        p = urlparse(href)
        if p.scheme in ("http", "https") and p.netloc.endswith("dzen.ru"):
            if p.path.startswith("/a/"):
                clean = p._replace(query="", fragment="")
                return clean.geturl()
        if href.startswith("/a/"):
            return f"https://dzen.ru{href.split('?')[0]}"
        return None
    except Exception as e:
        print(f"normalize error: {e}")
        return None


def scroll_page(driver: WebDriver, count: int, delay_seconds: int) -> None:
    for i in range(count):
        driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
        time.sleep(delay_seconds)
        if settings.DEBUG_LOG_SELECTORS:
            print(f"parser.debug Scrolling down... {i+1}/{count}")


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

# intranet/dzen_url_parser/parser.py
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

def parse_dzen_for_links_with_category(driver: WebDriver, link_storer: Optional[LinkStorer] = None) -> List[Dict[str, Optional[str]]]:
    print("Starting Dzen parsing process for all configured URLs...")
    unique: Dict[str, Optional[str]] = {}
    published_seen: Set[str] = set()

    try:
        for topic_url in settings.DZEN_ARTICLES_URLs:
            category = extract_category_from_topic_url(topic_url)

            driver.get(topic_url)
            # Явно ждём первый появившийся элемент карточки
            WebDriverWait(driver, 20).until(
                EC.any_of(
                    EC.presence_of_element_located((By.CSS_SELECTOR, 'a[data-testid="card-article-link"]')),
                    EC.presence_of_element_located((By.CSS_SELECTOR, 'a[data-testid="card-article-title-link"]'))
                )
            )

            # Скролл для догрузки
            for _ in range(settings.DZEN_SCROLL_COUNT):
                driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
                time.sleep(settings.DZEN_SCROLL_DELAY_SECONDS)

            urls = collect_links_from_html(driver.page_source)

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
        print(f"parser error during parsing: {e}")
        return [{"url": u, "category": c} for u, c in unique.items()]

def collect_links_from_html(html: str) -> List[str]:
    soup = BeautifulSoup(html, "html.parser")
    found: List[str] = []

    a_candidates = []
    for dtid in ("card-article-link", "card-article-title-link"):
        a_candidates.extend(soup.find_all("a", attrs={"data-testid": dtid}))

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

    # Расширенный фолбэк по href
    if not a_candidates:
        a_candidates = soup.select('a[href^="https://dzen.ru/a/"], a[href^="/a/"]')

    for a in a_candidates:
        href = a.get("href")
        if not href:
            continue
        norm = normalize_dzen_article_url(href)
        if norm:
            found.append(norm)

    seen: Set[str] = set()
    uniq: List[str] = []
    for url in found:
        if url not in seen:
            uniq.append(url)
            seen.add(url)

    if settings.DEBUG_LOG_SELECTORS:
        preview = ", ".join(uniq[:3])
        print(f"parser.debug collected={len(uniq)} preview={preview}")

    return uniq
