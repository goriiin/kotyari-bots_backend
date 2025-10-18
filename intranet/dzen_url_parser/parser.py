from __future__ import annotations

import time
from typing import Set, List
from urllib.parse import urlparse, urlunparse

from bs4 import BeautifulSoup
from selenium.webdriver.remote.webdriver import WebDriver

from .config import settings


def _normalize_dzen_article_url(href: str) -> str | None:
    """
    Нормализует ссылки статей Дзен к виду:
      https://dzen.ru/a/<id>
    Отбрасывает query/fragment, пропускает не-article ссылки.
    """
    try:
        p = urlparse(href)
        # Разрешаем абсолютные ссылки на dzen.ru/a/<...>
        if p.scheme in ("http", "https") and p.netloc.endswith("dzen.ru"):
            # Прямые статьи
            if p.path.startswith("/a/"):
                clean = p._replace(query="", fragment="")
                return urlunparse(clean)

        # Витринные ссылки иногда попадают относительными, проверим такой случай
        if href.startswith("/a/"):
            return f"https://dzen.ru{href.split('?')[0]}"

        # Игнорируем away?to=... (внешние ссылки) и иные разделы
        return None
    except Exception:
        return None


def _collect_links_from_html(html: str) -> List[str]:
    """
    Извлекает ссылки на статьи из HTML:
    1) Сначала — по ожидаемым data-testid у карточек.
    2) Fallback — любые якоря вида https(s)://dzen.ru/a/<...> и относительные /a/<...>.
    """
    soup = BeautifulSoup(html, "html.parser")

    found: List[str] = []

    # Карточки витрины — разные варианты data-testid на ленте/топиках
    a_candidates = []
    for dtid in ("card-article-link", "card-article-title-link"):
        a_candidates.extend(soup.find_all("a", attrs={"data-testid": dtid}))

    # Иногда карточка — это <article data-testid="floor-image-card"> с <a data-testid="card-article-link">
    for article in soup.find_all(attrs={"data-testid": "floor-image-card"}):
        a = article.find("a", attrs={"data-testid": "card-article-link"})
        if a is not None:
            a_candidates.append(a)

    # Дополнительные карточки в адаптивной сетке топиков
    for wrap in soup.find_all(attrs={"data-testid": "card-part-title"}):
        parent = wrap.parent
        if parent:
            a = parent.find("a", attrs={"data-testid": "card-article-title-link"})
            if a is not None:
                a_candidates.append(a)

    # Сбор href из карточек
    for a in a_candidates:
        href = a.get("href")
        if not href:
            continue
        norm = _normalize_dzen_article_url(href)
        if norm:
            found.append(norm)

    # Fallback — любые <a href="...">, совпадающие с шаблоном статей
    if not found:
        for a in soup.find_all("a", href=True):
            norm = _normalize_dzen_article_url(a["href"])
            if norm:
                found.append(norm)

    # Дедуп в порядке появления
    seen: Set[str] = set()
    uniq: List[str] = []
    for url in found:
        if url not in seen:
            uniq.append(url)
            seen.add(url)

    if settings.DEBUG_LOG_SELECTORS:
        # Печатаем первые ссылки для диагностики
        preview = ", ".join(uniq[:3])
        print(f"[parser.debug] collected={len(uniq)} preview=[{preview}]")

    return uniq


def _scroll_page(driver: WebDriver, count: int, delay_seconds: int) -> None:
    """
    Скроллит страницу вниз count раз с задержкой, чтобы подгрузить карточки.
    """
    for i in range(count):
        driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
        time.sleep(delay_seconds)
        if settings.DEBUG_LOG_SELECTORS:
            print(f"[parser.debug] Scrolling down... ({i+1}/{count})")


def _extract_category_from_topic_url(u: str) -> str | None:
    p = urlparse(u)
    path = (p.path or "/").rstrip("/")
    if not path or path == "/":
        return None
    return path.split("/")[-1]  # 'travel' для /topic/travel, 'articles' для /articles


def parse_dzen_for_links_with_category(driver: WebDriver) -> list[dict]:
    """
    Парсит набор витрин/топиков из настроек, скроллит и собирает ссылки статей,
    возвращая список словарей {"url": ..., "category": ...}.
    ВНИМАНИЕ: управление жизненным циклом driver оставлено вызывающему (gRPC слою).
    """
    print("Starting Dzen parsing process for all configured URLs...")
    unique: dict[str, str] = {}  # url -> category (первая встретившаяся)
    try:
        for topic_url in settings.DZEN_ARTICLES_URLs:
            category = _extract_category_from_topic_url(topic_url)
            driver.get(topic_url)
            time.sleep(5)  # базовая стабилизация после загрузки
            _scroll_page(driver, settings.DZEN_SCROLL_COUNT, settings.DZEN_SCROLL_DELAY_SECONDS)
            urls = _collect_links_from_html(driver.page_source)
            for u in urls:
                # сохраняем первую встретившуюся категорию
                if u not in unique:
                    unique[u] = category
        items = [{"url": u, "category": c} for u, c in unique.items()]
        return items
    except Exception as e:
        print(f"[parser] error during parsing: {e}")
        return []
