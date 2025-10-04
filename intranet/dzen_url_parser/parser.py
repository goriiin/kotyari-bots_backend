import time
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from urllib.parse import urlparse

from config import settings

def _get_webdriver() -> webdriver.Chrome:
    """Configures and returns a headless Chrome WebDriver instance."""
    options = webdriver.ChromeOptions()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--log-level=3")
    options.add_argument(
        "user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    )
    options.add_experimental_option("prefs", {"profile.managed_default_content_settings.images": 2})

    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=options)

def parse_dzen_for_links() -> list[str]:
    """
    Iterates through a list of Dzen topic URLs, scrolls each one to load
    dynamic content, and parses them to extract unique article links.

    Returns:
        A list of all unique URL strings found across all pages.
    """
    print("Starting Dzen parsing process for all configured URLs...")
    driver = None
    all_unique_links = set()

    try:
        driver = _get_webdriver()

        for topic_url in settings.DZEN_ARTICLES_URLs:
            print("-" * 50)
            print(f"Navigating to topic: {topic_url}")

            driver.get(topic_url)
            time.sleep(5)

            for i in range(settings.DZEN_SCROLL_COUNT):
                print(f"Scrolling down... ({i + 1}/{settings.DZEN_SCROLL_COUNT})")
                driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
                time.sleep(settings.DZEN_SCROLL_DELAY_SECONDS)

            # Parse the page source
            soup = BeautifulSoup(driver.page_source, 'html.parser')
            article_cards = soup.find_all('div', attrs={'data-testid': 'article-showcase-card'})

            found_on_page = 0
            for card in article_cards:
                link_tag = card.find('a', attrs={'data-testid': 'card-article-link'})
                if link_tag and link_tag.has_attr('href'):
                    href = link_tag['href']
                    parsed_url = urlparse(href)
                    clean_link = f"{parsed_url.scheme}://{parsed_url.netloc}{parsed_url.path}"

                    if clean_link not in all_unique_links:
                        all_unique_links.add(clean_link)
                        found_on_page += 1

            print(f"Found {found_on_page} new unique links on this page.")

        print("-" * 50)
        print(f"Parsing complete. Found a total of {len(all_unique_links)} unique article links.")
        return list(all_unique_links)

    except Exception as e:
        print(f"An error occurred during parsing: {e}")
        return []
    finally:
        if driver:
            print("Closing WebDriver.")
            driver.quit()