from __future__ import annotations

import os
from typing import Optional

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service


def create_anti_detect_driver(proxy: Optional[str] = None) -> webdriver.Chrome:
    """
    Plain Selenium Chrome driver without selenium-wire MITM proxy.
    Accepts proxy in formats:
      - "host port user pass"
      - "host port"
    Honors CHROME_BIN and CHROMEDRIVER_BIN if present.
    """
    opts = Options()
    opts.add_argument("--headless=new")
    opts.add_argument("--no-sandbox")
    opts.add_argument("--disable-dev-shm-usage")
    opts.add_argument("--disable-gpu")
    opts.add_argument("--disable-blink-features=AutomationControlled")
    opts.add_experimental_option("excludeSwitches", ["enable-automation"])
    opts.add_experimental_option("useAutomationExtension", False)

    chrome_bin = os.environ.get("CHROME_BIN")
    if chrome_bin:
        opts.binary_location = chrome_bin

    if proxy:
        parts = proxy.split()
        if len(parts) == 4:
            host, port, username, password = parts
            opts.add_argument(f"--proxy-server=http://{username}:{password}@{host}:{port}")
        elif len(parts) == 2:
            host, port = parts
            opts.add_argument(f"--proxy-server=http://{host}:{port}")

    chromedriver_bin = os.environ.get("CHROMEDRIVER_BIN")
    service = Service(chromedriver_bin) if chromedriver_bin else Service()

    driver = webdriver.Chrome(service=service, options=opts)

    # Stealth: navigator.webdriver = undefined
    driver.execute_cdp_cmd(
        "Page.addScriptToEvaluateOnNewDocument",
        {
            "source": """
                Object.defineProperty(navigator, 'webdriver', {
                  get: () => undefined
                });
            """
        },
    )

    return driver
