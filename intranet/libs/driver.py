import logging
import os
from typing import Optional, Dict, Any

from seleniumwire import webdriver  # ensures proxy auth works reliably
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.chrome.options import Options as ChromeOptions

logger = logging.getLogger(__name__)

def service() -> ChromeService:
    candidates = [os.environ.get("CHROMEDRIVERBIN"), "/usr/bin/chromedriver"]
    for path in candidates:
        if path and os.path.exists(path):
            logger.info(f"Using chromedriver path: {path}")
            return ChromeService(executable_path=path)
    logger.info("Using auto-detected chromedriver")
    return ChromeService()

def parseproxy(proxy: str) -> Optional[Dict[str, str]]:
    """
    Accepts 'ip port user pass' or 'ip:port:user:pass'.
    Returns a Selenium Wire proxy dict ready for seleniumwire_options['proxy'].
    """
    try:
        # Normalize separators (spaces or colons)
        parts = proxy.replace(":", " ").split()
        if len(parts) != 4:
            raise ValueError("Proxy format must be 'ip port user pass' or 'ip:port:user:pass'")
        host, port, username, password = parts
        http_url = f"http://{username}:{password}@{host}:{port}"
        https_url = f"https://{username}:{password}@{host}:{port}"
        return {
            "http": http_url,
            "https": https_url,
            "no_proxy": "localhost,127.0.0.1",
        }
    except Exception as e:
        logger.error(f"Proxy parse error: {e} — running without proxy")
        return None

def build_chrome_options() -> ChromeOptions:
    opts = ChromeOptions()
    opts.add_argument("--disable-blink-features=AutomationControlled")
    opts.add_argument("--no-sandbox")
    opts.add_argument("--disable-dev-shm-usage")
    opts.add_argument("--lang=ru-RU")
    opts.add_experimental_option("excludeSwitches", ["enable-automation"])
    opts.add_experimental_option("useAutomationExtension", False)
    if os.getenv("CHROME_HEADLESS", "1") in ("1", "true", "True"):
        opts.add_argument("--headless=new")
    return opts

def createantidetectdriver(proxy: Optional[str] = None) -> webdriver.Chrome:
    page_load_timeout = int(os.getenv("SELENIUM_PAGELOAD_TIMEOUT", "180"))
    implicit_wait = int(os.getenv("SELENIUM_IMPLICIT_WAIT", "10"))
    conn_timeout = int(os.getenv("SELENIUM_CONNECTION_TIMEOUT", "60"))
    req_timeout = int(os.getenv("SELENIUM_REQUEST_TIMEOUT", "180"))

    chrome_opts = build_chrome_options()

    # IMPORTANT: widen type to allow nested dict in 'proxy'
    seleniumwire_options: Dict[str, Any] = {
        "suppress_connection_errors": True,
        "connection_timeout": conn_timeout,
        "request_timeout": req_timeout,
    }

    if proxy:
        proxy_cfg = parseproxy(proxy)
        if proxy_cfg:
            # Selenium Wire expects a dict here: {'http': '...', 'https': '...', 'no_proxy': '...'}
            seleniumwire_options["proxy"] = proxy_cfg

    drv = webdriver.Chrome(service=service(), options=chrome_opts, seleniumwire_options=seleniumwire_options)
    try:
        drv.set_page_load_timeout(page_load_timeout)
    except Exception as e:
        logger.error(f"Failed to set page load timeout in driver: {e}")
        pass
    try:
        drv.implicitly_wait(implicit_wait)
    except Exception as e:
        logger.error(f"Failed to set implicit wait timeout in driver: {e}")
        pass

    try:
        drv.execute_cdp_cmd("Emulation.setTimezoneOverride", {"timezoneId": os.getenv("TIMEZONE", "Europe/Moscow")})
    except Exception as e:
        logger.warning(f"Failed to set timezone override: {e}")

    logger.info("Anti-detect driver created")
    return drv
