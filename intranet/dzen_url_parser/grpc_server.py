from __future__ import annotations

import logging

import grpc
from google.protobuf import empty_pb2
from api.protos.url_fetcher.gen import start_fetching_pb2_grpc
from .parser import parse_dzen_for_links_with_category
from .redis_adapter import LinkStorer, RedisPublisherAdapter
from intranet.libs.driver import createantidetectdriver
from .config import settings
from ..libs.proxy_pool import ProxyPool

logger = logging.getLogger(__name__)

class ProfileServiceServicer(start_fetching_pb2_grpc.ProfileServiceServicer):
    def __init__(self, linkstorer: LinkStorer | None = None):
        print("ProfileServiceServicer initialized.")
        self.linkstorer = linkstorer or RedisPublisherAdapter()
        self.proxypool = ProxyPool(settings.PROXY_FILEPATH)

    def StartFetching(self, request, context):
        print("gRPC call received: StartFetching.")
        driver = None
        try:
            proxy = self.proxypool.get_random_proxy()
            print(f"Selected proxy: {proxy.split()[0] if proxy else None}")
            driver = createantidetectdriver(proxy=proxy)
            links = parse_dzen_for_links_with_category(driver, linkstorer=self.linkstorer)
            print(f"Adapter will publish {len(links)} links...")
            published = self.linkstorer.store_links(links)
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
                except Exception as e:
                    logger.error(f"Failed to close driver: {e}")
                    pass
