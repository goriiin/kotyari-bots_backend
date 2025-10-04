from intranet.libs.redis import RedisClient
from config import settings
from abc import ABC, abstractmethod

class LinkStorer(ABC):
    """
    Defines the interface for any class that stores parsed links.
    This is the abstraction that our high-level components will depend on.
    """

    @abstractmethod
    def store_links(self, links: list[str]) -> int:
        """
        Stores a list of links in a data store.

        Args:
            links: A list of URL strings.

        Returns:
            The number of new links successfully stored.

        Raises:
            ConnectionError: If the data store is unavailable.
        """
        pass

class RedisPublisherAdapter(LinkStorer):
    """
    An adapter that implements the LinkStorerInterface but uses the new
    RedisClient publisher internally. This decouples our servicer from the
    specifics of the new Redis library.
    """
    def __init__(self):
        """
        Initializes the underlying RedisClient.
        """
        try:
            self.redis_client = RedisClient(
                host=settings.REDIS_HOST,
                port=settings.REDIS_PORT,
                processed_urls_key=settings.REDIS_PROCESSED_URLS_KEY
            )
        except Exception as e:
            # If the client fails to connect, we can't do anything.
            print(f"Fatal: Failed to initialize RedisPublisherAdapter: {e}")
            self.redis_client = None

    def store_links(self, links: list[str]) -> int:
        """
        Iterates through links and publishes them one by one using the
        new client's publish method.
        """
        if not self.redis_client:
            raise ConnectionError("Redis client in adapter is not available.")

        new_links_published = 0
        for link in links:
            # The publish method returns True if the URL was new and published.
            if self.redis_client.publish(topic=settings.REDIS_PUBLISH_TOPIC, url=link):
                new_links_published += 1

        print(f"Adapter finished. Published {new_links_published} new links to topic '{settings.REDIS_PUBLISH_TOPIC}'.")
        return new_links_published