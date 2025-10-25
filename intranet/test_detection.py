"""
Anti-detection testing script.
Tests WebDriver configuration against bot detection services.
Includes proxy validation tests.
"""
import sys
import logging
import requests
from libs.driver import create_anti_detect_driver
from libs.proxy_pool import ProxyPool

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

logger = logging.getLogger(__name__)


def test_proxy_requests(proxy_str: str) -> bool:
    """Test proxy connectivity using requests library."""
    try:
        parts = proxy_str.split(":")
        if len(parts) != 4:
            logger.error("Invalid proxy format")
            return False

        host, port, username, password = parts
        proxy_url = f"http://{username}:{password}@{host}:{port}/"

        proxies = {
            "http": proxy_url,
            "https": proxy_url
        }

        logger.info(f"Testing proxy connectivity: {host}:{port}")
        response = requests.get(
            "https://ipv4.webshare.io/",
            proxies=proxies,
            timeout=10
        )

        logger.info(f"Proxy response: {response.text.strip()}")
        logger.info(f"Status code: {response.status_code}")

        if response.status_code == 200:
            logger.info("✓ Proxy works with requests library")
            return True
        else:
            logger.error(f"✗ Proxy returned status {response.status_code}")
            return False

    except requests.exceptions.ProxyError as e:
        logger.error(f"✗ Proxy connection error: {e}")
        return False
    except requests.exceptions.Timeout:
        logger.error("✗ Proxy timeout")
        return False
    except Exception as e:
        logger.error(f"✗ Unexpected error: {e}")
        return False


def main():
    print("=" * 60)
    print("ANTI-DETECTION & PROXY TESTING")
    print("=" * 60)

    # Test 1: Proxy pool loading
    print("\n[1/5] Testing Proxy Pool...")
    pool = ProxyPool("ip.txt")
    proxy = pool.get_random_proxy()

    if not proxy:
        print("❌ No proxies loaded. Check file path.")
        sys.exit(1)

    print(f"✓ Loaded {len(pool.proxies)} proxies")
    print(f"✓ Selected proxy: {proxy.split(':')[0]}:{proxy.split(':')[1]}")

    # Test 2: Proxy connectivity with requests
    print("\n[2/5] Testing proxy with requests library...")
    proxy_works = test_proxy_requests(proxy)

    if not proxy_works:
        print("⚠ Proxy test failed, but continuing with Selenium tests...")
        print("  This might be a connectivity issue or proxy service problem.")

    # Test 3: Driver creation
    print("\n[3/5] Creating anti-detect driver with proxy...")
    try:
        driver = create_anti_detect_driver(proxy=proxy)
        print("✓ Driver created successfully")
    except Exception as e:
        print(f"❌ Failed to create driver: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

    # Test 4: IP check via proxy
    print("\n[4/5] Checking IP through Selenium + Proxy...")
    try:
        driver.get("https://httpbin.org/ip")
        import time
        time.sleep(2)
        page_text = driver.find_element("tag name", "body").text
        print(f"Response from httpbin.org/ip:\n{page_text}")

        # Verify IP matches proxy
        expected_ip = proxy.split(':')[0]
        if expected_ip in page_text:
            print(f"✓ IP matches proxy: {expected_ip}")
        else:
            print(f"⚠ IP doesn't match proxy. Expected {expected_ip}")
            print("  Extension might not be loaded or proxy authentication failed.")
    except Exception as e:
        print(f"⚠ Error checking IP: {e}")

    # Test 5: Bot detection check
    print("\n[5/5] Testing bot detection (bot.sannysoft.com)...")
    try:
        driver.get("https://bot.sannysoft.com/")
        import time
        time.sleep(3)

        page_source = driver.page_source.lower()

        print("\n" + "=" * 60)
        print("BOT DETECTION RESULTS:")

        if "webdriver" in page_source:
            if "missing" in page_source or "false" in page_source:
                print("  ✓ navigator.webdriver: PASSED (missing/false)")
            else:
                print("  ✗ navigator.webdriver: DETECTED")

        if "chrome" in page_source:
            print("  ✓ Chrome detected")

        print("=" * 60)

    except Exception as e:
        print(f"⚠ Error loading detection page: {e}")

    # Cleanup
    driver.quit()
    print("\n✓ All tests completed.")
    print("\nSUMMARY:")
    print(f"  - Proxy pool: {len(pool.proxies)} proxies loaded")
    print(f"  - Selected proxy: {proxy.split(':')[0]}:{proxy.split(':')[1]}")
    print(f"  - Proxy connectivity: {'✓ PASS' if proxy_works else '✗ FAIL'}")
    print("  - Selenium driver: ✓ PASS")
    print("  - Bot detection: Check results above")


if __name__ == "__main__":
    main()
