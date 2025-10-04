import time
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager
from urllib.parse import urlparse

# Import our centralized settings
from config import settings

def _get_webdriver() -> webdriver.Chrome:
    """Configures and returns a headless Chrome WebDriver instance."""
    options = webdriver.ChromeOptions()
    options.add_argument("--headless")  # Run in the background without a UI
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--log-level=3") # Suppress unnecessary console logs
    options.add_argument(
        "user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    )
    # Disable image loading to speed up page loads
    options.add_experimental_option("prefs", {"profile.managed_default_content_settings.images": 2})

    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=options)

def parse_dzen_for_links() -> list[str]:
    """
    Fetches the Dzen articles page, scrolls to load dynamic content,
    and parses it to extract unique article links.

    Returns:
        A list of unique URL strings found on the page.
    """
    print("Starting Dzen parsing process...")
    driver = None
    try:
        driver = _get_webdriver()
        unique_links = set()
        for url in settings.DZEN_ARTICLES_URLs:
            print(f"Navigating to {url}")
            driver.get(settings.DZEN_ARTICLES_URL)
            time.sleep(5) # Wait for the initial page to render

            # Scroll down the page multiple times to trigger dynamic content loading
            for i in range(settings.DZEN_SCROLL_COUNT):
                print(f"Scrolling down... ({i + 1}/{settings.DZEN_SCROLL_COUNT})")
                driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
                time.sleep(settings.DZEN_SCROLL_DELAY_SECONDS)

            # Parse the final page source with Beautiful Soup
            soup = BeautifulSoup(driver.page_source, 'html.parser')

            # The selector targets the <a> tag inside each article card.
            # This is based on the data-testid attributes which are stable selectors.
            # NOTE: The selector `a[data-testid="card-article-link"]` is more direct.
            article_cards = soup.find_all('div', attrs={'data-testid': 'article-showcase-card'})


            for card in article_cards:
                link_tag = card.find('a', attrs={'data-testid': 'card-article-link'})
                if link_tag and link_tag.has_attr('href'):
                    href = link_tag['href']
                    # Clean the URL to remove tracking parameters and ensure consistency.
                    # We rebuild the URL with only the scheme, netloc, and path.
                    parsed_url = urlparse(href)
                    clean_link = f"{parsed_url.scheme}://{parsed_url.netloc}{parsed_url.path}"
                    unique_links.add(clean_link)

        print(f"Parsing complete. Found {len(unique_links)} unique article links.")
        return list(unique_links)

    except Exception as e:
        print(f"An error occurred during parsing: {e}")
        return [] # Return an empty list on failure to prevent crashes
    finally:
        if driver:
            print("Closing WebDriver.")
            driver.quit()