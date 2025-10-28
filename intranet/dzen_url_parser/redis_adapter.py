from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Iterable, Optional, TypedDict
import json

from intranet.libs.redis import RedisClient
from .config import settings


class LinkItem(TypedDict):
    url: str
    category: Optional[str]


class LinkStorer(ABC):
    @abstractmethod
    def store_links(self, links: Iterable[LinkItem]) -> int:
        raise NotImplementedError


class RedisPublisherAdapter(LinkStorer):
    """
    Publishes links to Redis PubSub and marks them in the processed ZSET.
    Expects items with shape: {"url": str, "category": Optional[str]}.
    """

    def __init__(self) -> None:
        self.topic = settings.REDIS_PUBLISH_TOPIC
        self.processed_key = settings.REDIS_PROCESSED_URLS_KEY

        username = settings.REDIS_USER or None
        password = settings.REDIS_PASSWORD

        # Keep argument names consistent with existing RedisClient initializer
        self.redisclient = RedisClient(
            host=settings.REDISHOST,
            port=settings.REDISPORT,
            processed_urls_key=self.processed_key,
            username=username,
            password=password,
        )

    def store_links(self, links: Iterable[LinkItem]) -> int:
        if not self.redisclient:
            raise ConnectionError("Redis client in adapter is not available.")

        new_published = 0
        for item in links:
            url = item["url"]
            category = item.get("category")
            payload = json.dumps(
                {"url": url, "category": category},
                ensure_ascii=False,
            )
            # Preserve the publish signature used elsewhere in the project
            if self.redisclient.publish(topic=self.topic, url=url, payload=payload):
                new_published += 1

        return new_published
