"""
Simple proxy connectivity test using requests library.
Quick validation without Selenium overhead.
"""
import sys
import requests
from libs.proxy_pool import ProxyPool


def test_proxy(proxy_str: str) -> dict:
    """Test proxy and return results."""
    try:
        parts = proxy_str.split(":")
        host, port, username, password = parts
        proxy_url = f"http://{username}:{password}@{host}:{port}/"

        proxies = {
            "http": proxy_url,
            "https": proxy_url
        }

        print(f"\nTesting proxy: {host}:{port}")
        print(f"Making request to https://ipv4.webshare.io/...")

        response = requests.get(
            "https://ipv4.webshare.io/",
            proxies=proxies,
            timeout=10
        )

        return {
            "success": True,
            "status_code": response.status_code,
            "ip": response.text.strip(),
            "proxy_ip": host
        }

    except Exception as e:
        return {
            "success": False,
            "error": str(e),
            "proxy_ip": proxy_str.split(':')[0]
        }


def main():
    print("=" * 60)
    print("SIMPLE PROXY CONNECTIVITY TEST")
    print("=" * 60)

    pool = ProxyPool("ip.txt")

    if not pool.proxies:
        print("❌ No proxies loaded")
        sys.exit(1)

    print(f"\nLoaded {len(pool.proxies)} proxies")
    print("Testing all proxies...\n")

    results = []
    for i, proxy in enumerate(pool.proxies, 1):
        print(f"[{i}/{len(pool.proxies)}] ", end="")
        result = test_proxy(proxy)
        results.append(result)

        if result["success"]:
            print(f"✓ PASS - IP: {result['ip']}")
        else:
            print(f"✗ FAIL - Error: {result['error']}")

    # Summary
    print("\n" + "=" * 60)
    print("SUMMARY")
    print("=" * 60)

    successful = sum(1 for r in results if r["success"])
    failed = len(results) - successful

    print(f"Total proxies: {len(results)}")
    print(f"Successful: {successful} ({successful / len(results) * 100:.1f}%)")
    print(f"Failed: {failed} ({failed / len(results) * 100:.1f}%)")

    if successful > 0:
        print("\nWorking proxies:")
        for r in results:
            if r["success"]:
                print(f"  - {r['proxy_ip']} → {r['ip']}")


if __name__ == "__main__":
    main()
