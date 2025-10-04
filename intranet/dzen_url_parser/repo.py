import redis
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

class RedisLinkStorer(LinkStorer):
    """
    A concrete implementation of the LinkStorerInterface that uses Redis
    as the data store.
    """
    def __init__(self):
        """
        Initializes the Redis client and establishes a connection.
        """
        self.redis_client = None
        try:
            self.redis_client = redis.Redis(
                host=settings.REDIS_HOST,
                port=settings.REDIS_PORT,
                db=0,
                decode_responses=True
            )
            self.redis_client.ping()
            print(f"RedisLinkStorer connected successfully to Redis at {settings.REDIS_HOST}:{settings.REDIS_PORT}")
        except redis.exceptions.ConnectionError as e:
            print(f"Fatal: RedisLinkStorer could not connect to Redis: {e}")
            # The client remains None, and subsequent calls will fail.

    def store_links(self, links: list[str]) -> int:
        """
        Stores a list of links in a Redis set.
        """
        if not self.redis_client:
            print("Cannot store links: Redis client is not available.")
            raise ConnectionError("Redis connection not established.")

        if not links:
            return 0

        # The sadd command returns the number of elements that were added to the set.
        num_added = self.redis_client.sadd("dzen_links", *links)
        print(f"Called Redis SADD for {len(links)} links. {num_added} were new.")
        return num_added