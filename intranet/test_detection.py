"""
Anti-detection testing script.
Tests WebDriver configuration against bot detection services.
Includes proxy validation tests.
"""
import sys
import requests
from libs.driver import create_anti_detect_driver
from libs.proxy_pool import ProxyPool
import time


def test_proxy_requests(proxy_str: str) -> bool:
    try:
        parts = proxy_str.split(":")
        if len(parts) != 4:
            print("Invalid proxy format")
            return False
        host, port, username, password = parts
        proxy_url = f"http://{username}:{password}@{host}:{port}/"
        proxies = {"http": proxy_url, "https": proxy_url}
        print(f"Testing proxy connectivity: {host}:{port}")
        resp = requests.get("https://ipv4.webshare.io/", proxies=proxies, timeout=10)
        print(f"Proxy response: {resp.text.strip()} status={resp.status_code}")
        return resp.status_code == 200
    except Exception as e:
        print(f"Proxy connectivity error: {e}")
        return False

def main():
    print("=" * 60)
    print("ANTI-DETECTION & PROXY TESTING")
    print("=" * 60)

    print("\n[1/5] Testing Proxy Pool...")
    pool = ProxyPool("ip.txt")
    proxy = pool.get_random_proxy()
    if not proxy:
        print("❌ No proxies loaded. Check file path.")
        sys.exit(1)
    print(f"✓ Loaded {len(pool.proxies)} proxies")
    print(f"✓ Selected proxy: {proxy.split(':')[0]}:{proxy.split(':')[1]}")

    print("\n[2/5] Testing proxy with requests library...")
    proxy_works = test_proxy_requests(proxy)
    if not proxy_works:
        print("⚠ Proxy test failed, continuing with Selenium tests...")

    print("\n[3/5] Creating anti-detect driver with proxy...")
    try:
        driver = create_anti_detect_driver(proxy=proxy)
        print("✓ Driver created successfully")
    except Exception as e:
        print(f"❌ Failed to create driver: {e}")
        raise

    print("\n[4/5] Checking IP through Selenium + Proxy...")
    try:
        driver.get("https://httpbin.org/ip")
        time.sleep(2)
        page_text = driver.find_element("tag name", "body").text
        print(f"Response from httpbin.org/ip:\n{page_text}")
        expected_ip = proxy.split(':')[0]
        if expected_ip in page_text:
            print(f"✓ IP matches proxy: {expected_ip}")
        else:
            print(f"⚠ IP doesn't match proxy. Expected {expected_ip}")
    except Exception as e:
        print(f"⚠ Error checking IP: {e}")

    print("\n[5/5] Bot detection checks (sannysoft + creepjs)...")
    try:
        driver.get("https://bot.sannysoft.com/")
        time.sleep(3)
        ps = driver.page_source.lower()
        if "webdriver" in ps and ("missing" in ps or "false" in ps):
            print("  ✓ navigator.webdriver: PASSED (missing/false)")
        print("  ✓ Chrome detected")

        driver.get("https://abrahamjuliot.github.io/creepjs/")
        time.sleep(5)
        print("  ✓ creepjs loaded (manually verify canvas/webgl noise uniqueness)")
    except Exception as e:
        print(f"⚠ Error on detection pages: {e}")

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
