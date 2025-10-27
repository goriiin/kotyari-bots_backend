"""
Anti-detect WebDriver factory with proxy support using selenium-wire.
"""
import os
import logging
import random
import tempfile
from typing import Optional

from seleniumwire import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.chrome.options import Options as ChromeOptions

logger = logging.getLogger(__name__)

# Realistic Linux Chrome UAs (rotate periodically)
USER_AGENTS = [
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.6261.128 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.6167.140 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.6099.224 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.6045.199 Safari/537.36",
]

# Pool of real GPU vendors/renderers (exclude SwiftShader)
GPU_VENDORS = [
    ("Intel Inc.", "Intel(R) UHD Graphics 620"),
    ("Intel Inc.", "Mesa Intel(R) UHD Graphics 630"),
    ("NVIDIA Corporation", "NVIDIA GeForce GTX 1050"),
    ("NVIDIA Corporation", "NVIDIA GeForce GTX 1650"),
    ("AMD", "Radeon RX 560 Series"),
]

LANGS = ["ru-RU", "ru", "en-US"]
TIMEZONE = "Europe/Moscow"


def _build_chrome_options() -> ChromeOptions:
    opts = ChromeOptions()

    # Binary override if provided
    chrome_bin = os.environ.get("CHROME_BIN")
    if chrome_bin and os.path.exists(chrome_bin):
        opts.binary_location = chrome_bin
        logger.info(f"Using Chrome binary: {chrome_bin}")

    # Headless and stability flags
    opts.add_argument("--headless=new")
    opts.add_argument("--no-sandbox")
    opts.add_argument("--disable-dev-shm-usage")
    opts.add_argument("--disable-gpu")
    opts.add_argument("--window-size=1920,1080")
    opts.add_argument("--lang=ru-RU")

    # Anti-automation flags
    opts.add_experimental_option("excludeSwitches", ["enable-automation"])
    opts.add_experimental_option("useAutomationExtension", False)
    opts.add_argument("--disable-blink-features=AutomationControlled")

    # Preferences
    opts.add_experimental_option(
        "prefs",
        {
            "intl.accept_languages": ",".join(LANGS),
            "profile.managed_default_content_settings.images": 2,
        },
    )
    return opts


def _apply_stealth(driver) -> None:
    """
    Apply stealth techniques via Chrome DevTools Protocol:
    - navigator.webdriver = undefined
    - WebGL vendor/renderer spoof
    - navigator.languages spoof
    - Canvas/WebGL noise (1–2%)
    - Timezone override to Europe/Moscow
    """
    vendor, renderer = random.choice(GPU_VENDORS)
    # Small deterministic seed per session to keep canvas stable per run
    seed = str(random.randint(10_000, 99_999))

    stealth_js = f"""
// navigator.webdriver -> undefined
Object.defineProperty(navigator, 'webdriver', {{get: () => undefined}});

// navigator.languages spoof
Object.defineProperty(navigator, 'languages', {{ get: () => {LANGS} }});

// WebGL vendor/renderer spoof
(function() {{
  const getParameter = WebGLRenderingContext.prototype.getParameter;
  WebGLRenderingContext.prototype.getParameter = function(parameter) {{
    if (parameter === 37445) return '{vendor}';
    if (parameter === 37446) return '{renderer}';
    return getParameter.apply(this, arguments);
  }};
  if (window.WebGL2RenderingContext) {{
    const getParameter2 = WebGL2RenderingContext.prototype.getParameter;
    WebGL2RenderingContext.prototype.getParameter = function(parameter) {{
      if (parameter === 37445) return '{vendor}';
      if (parameter === 37446) return '{renderer}';
      return getParameter2.apply(this, arguments);
    }};
  }}
}})();

// Canvas noise (1–2%)
(function() {{
  const SEED = '{seed}';
  function hash(x, y, c) {{
    let h = 2166136261;
    for (const ch of (SEED + ':' + x + ':' + y + ':' + c)) {{
      h ^= ch.charCodeAt(0);
      h = Math.imul(h, 16777619);
    }}
    return (h >>> 0) / 4294967295;
  }}
  const _getImageData = CanvasRenderingContext2D.prototype.getImageData;
  CanvasRenderingContext2D.prototype.getImageData = function(sx, sy, sw, sh) {{
    const data = _getImageData.apply(this, arguments);
    try {{
      const arr = data.data;
      for (let y = 0; y < sh; y++) {{
        for (let x = 0; x < sw; x++) {{
          const i = (y * sw + x) * 4;
          for (let c = 0; c < 3; c++) {{
            const rnd = hash(sx + x, sy + y, c);
            const delta = (rnd < 0.5 ? -1 : 1) * (1 + Math.floor(rnd * 2)); // 1–2
            arr[i + c] = Math.max(0, Math.min(255, arr[i + c] + delta));
          }}
        }}
      }}
    }} catch (e) {{}}
    return data;
  }};
  const _toDataURL = HTMLCanvasElement.prototype.toDataURL;
  HTMLCanvasElement.prototype.toDataURL = function() {{
    try {{
      const ctx = this.getContext('2d');
      if (ctx) {{
        const img = ctx.getImageData(0, 0, this.width, this.height);
        ctx.putImageData(img, 0, 0);
      }}
    }} catch (e) {{}}
    return _toDataURL.apply(this, arguments);
  }};
}})();
"""
    try:
        driver.execute_cdp_cmd("Page.addScriptToEvaluateOnNewDocument", {"source": stealth_js})
    except Exception as e:
        logger.warning(f"Failed to inject stealth JS: {e}")

    # Timezone override
    try:
        driver.execute_cdp_cmd("Emulation.setTimezoneOverride", {"timezoneId": TIMEZONE})
    except Exception as e:
        logger.warning(f"Failed to set timezone override: {e}")


def _service() -> ChromeService:
    candidates = [
        os.environ.get("CHROMEDRIVER_BIN"),
        "/usr/bin/chromedriver",
    ]
    for path in candidates:
        if path and os.path.exists(path):
            logger.info(f"Using chromedriver: {path}")
            return ChromeService(executable_path=path)
    logger.info("Using auto-detected chromedriver")
    return ChromeService()


def _parse_proxy(proxy: str) -> Optional[dict]:
    try:
        parts = proxy.split(":")
        if len(parts) != 4:
            raise ValueError("Proxy format must be ip:port:user:pass")
        host, port, username, password = parts
        return {
            "http": f"http://{username}:{password}@{host}:{port}",
            "https": f"http://{username}:{password}@{host}:{port}",
            "no_proxy": "localhost,127.0.0.1",
        }
    except Exception as e:
        logger.error(f"Proxy parse error: {e}; running without proxy")
        return None


def create_anti_detect_driver(proxy: Optional[str] = None) -> webdriver.Chrome:
    """
    Create Chrome WebDriver with anti-detection and optional proxy auth.
    Uses selenium-wire for reliable proxy authentication.
    """
    chrome_opts = _build_chrome_options()

    # Random realistic UA
    ua = random.choice(USER_AGENTS)
    chrome_opts.add_argument(f"user-agent={ua}")
    logger.info(f"Using User-Agent: {ua}")

    # selenium-wire proxy config
    seleniumwire_options = {}
    if proxy:
        parsed = _parse_proxy(proxy)
        if parsed:
            seleniumwire_options = {"proxy": parsed}
            try:
                # Mask the proxy auth header leakage to target origins
                seleniumwire_options["suppress_connection_errors"] = False
            except Exception:
                pass
            host, port = proxy.split(":")[0], proxy.split(":")[1]
            logger.info(f"Proxy configured: {host}:{port}")

    driver = webdriver.Chrome(
        service=_service(),
        options=chrome_opts,
        seleniumwire_options=seleniumwire_options or None,
    )

    # Apply stealth after driver creation
    _apply_stealth(driver)

    logger.info("Anti-detect driver created")
    return driver
