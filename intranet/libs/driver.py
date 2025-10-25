"""
Anti-detect WebDriver factory with proxy support using selenium-wire.
"""
import os
import logging
import random
from typing import Optional

from seleniumwire import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.chrome.options import Options as ChromeOptions

logger = logging.getLogger(__name__)

USER_AGENTS = [
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
]

GPU_VENDORS = [
    ("Intel Inc.", "Intel(R) HD Graphics 620"),
    ("NVIDIA Corporation", "NVIDIA GeForce GTX 1050"),
]


def create_anti_detect_driver(proxy: Optional[str] = None) -> webdriver.Chrome:
    """
    Creates Chrome WebDriver with anti-detection and proxy support.
    Uses selenium-wire for reliable proxy authentication.
    """
    chrome_opts = ChromeOptions()

    # Chrome binary
    chrome_bin = os.environ.get("CHROME_BIN")
    if chrome_bin and os.path.exists(chrome_bin):
        chrome_opts.binary_location = chrome_bin
        logger.info(f"Using Chrome binary: {chrome_bin}")

    # Headless mode
    chrome_opts.add_argument("--headless=new")
    chrome_opts.add_argument("--no-sandbox")
    chrome_opts.add_argument("--disable-dev-shm-usage")
    chrome_opts.add_argument("--disable-gpu")
    chrome_opts.add_argument("--window-size=1920,1080")

    # User-Agent
    user_agent = random.choice(USER_AGENTS)
    chrome_opts.add_argument(f"user-agent={user_agent}")
    logger.info(f"Using User-Agent: {user_agent}")

    # Languages
    chrome_opts.add_argument("--lang=ru-RU")
    chrome_opts.add_experimental_option("prefs", {
        "intl.accept_languages": "ru-RU,ru,en-US",
        "profile.managed_default_content_settings.images": 2,
    })

    # Anti-automation
    chrome_opts.add_experimental_option("excludeSwitches", ["enable-automation"])
    chrome_opts.add_experimental_option("useAutomationExtension", False)
    chrome_opts.add_argument("--disable-blink-features=AutomationControlled")

    # Selenium-wire options with proxy
    seleniumwire_options = {}

    if proxy:
        parts = proxy.split(":")
        if len(parts) == 4:
            host, port, username, password = parts
            seleniumwire_options = {
                'proxy': {
                    'http': f'http://{username}:{password}@{host}:{port}',
                    'https': f'http://{username}:{password}@{host}:{port}',
                    'no_proxy': 'localhost,127.0.0.1'
                }
            }
            logger.info(f"Proxy configured: {host}:{port}")

    # ChromeDriver service
    chromedriver_paths = [
        os.environ.get("CHROMEDRIVER_BIN"),
        "/usr/bin/chromedriver",
    ]

    chromedriver_path = None
    for path in chromedriver_paths:
        if path and os.path.exists(path):
            chromedriver_path = path
            break

    if chromedriver_path:
        service = ChromeService(executable_path=chromedriver_path)
        logger.info(f"Using chromedriver: {chromedriver_path}")
    else:
        service = ChromeService()
        logger.info("Using auto-detected chromedriver")

    # Create driver with selenium-wire
    driver = webdriver.Chrome(
        service=service,
        options=chrome_opts,
        seleniumwire_options=seleniumwire_options
    )

    # Apply stealth scripts
    _apply_stealth_scripts(driver)

    logger.info("Anti-detect driver created")
    return driver


def _apply_stealth_scripts(driver) -> None:
    """Apply stealth via CDP."""
    vendor, renderer = random.choice(GPU_VENDORS)

    stealth_js = f"""
    Object.defineProperty(navigator, 'webdriver', {{get: () => undefined}});
    const getParameter = WebGLRenderingContext.prototype.getParameter;
    WebGLRenderingContext.prototype.getParameter = function(parameter) {{
        if (parameter === 37445) return '{vendor}';
        if (parameter === 37446) return '{renderer}';
        return getParameter.apply(this, arguments);
    }};
    """

    try:
        driver.execute_cdp_cmd("Page.addScriptToEvaluateOnNewDocument", {"source": stealth_js})
        logger.debug("Stealth scripts applied")
    except Exception as e:
        logger.warning(f"Failed to apply stealth: {e}")
