from __future__ import annotations

import grpc
from google.protobuf import empty_pb2

from api.protos.url_fetcher.gen import start_fetching_pb2_grpc
# ВАЖНО: импортируем функцию напрямую, чтобы упасть на этапе импорта, если её нет
from .parser import parse_dzen_for_links_with_category
from .redis_adapter import LinkStorer, RedisPublisherAdapter

from selenium import webdriver
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.chrome.options import Options as ChromeOptions


def _make_headless_chrome() -> webdriver.Chrome:
    chrome_opts = ChromeOptions()
    chrome_opts.add_argument("--headless=new")
    chrome_opts.add_argument("--no-sandbox")
    chrome_opts.add_argument("--disable-dev-shm-usage")
    chrome_opts.add_argument("--disable-gpu")
    chrome_opts.add_argument("--window-size=1920,1080")
    chrome_opts.add_argument("--disable-extensions")
    chrome_opts.add_argument("--disable-infobars")
    chrome_opts.add_argument("--lang=ru_RU.UTF-8")
    service = ChromeService(executable_path="/usr/bin/chromedriver")
    return webdriver.Chrome(service=service, options=chrome_opts)


class ProfileServiceServicer(start_fetching_pb2_grpc.ProfileServiceServicer):
    def __init__(self, link_storer: LinkStorer) -> None:
        print("ProfileServiceServicer initialized.")
        self.link_storer = link_storer or RedisPublisherAdapter()

    def StartFetching(self, request, context):
        print("gRPC call received: StartFetching.")
        driver = None
        try:
            driver = _make_headless_chrome()
            links = parse_dzen_for_links_with_category(driver)
            print(f"Adapter will publish {len(links)} links...")
            published = self.link_storer.store_links( links)
            print(f"Adapter finished. Published {published} links.")
            return empty_pb2.Empty()
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            context.set_details("An internal server error occurred during fetching.")
            context.set_code(grpc.StatusCode.INTERNAL)
            return empty_pb2.Empty()
        finally:
            if driver is not None:
                try:
                    driver.quit()
                    print("WebDriver closed.")
                except Exception:
                    pass
