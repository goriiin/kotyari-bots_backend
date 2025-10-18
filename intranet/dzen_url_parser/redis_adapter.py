from __future__ import annotations

from abc import ABC, abstractmethod
from typing import Iterable

from intranet.libs.redis import RedisClient  # предполагаемый существующий клиент
from .config import settings

import json

class LinkStorer(ABC):
    @abstractmethod
    def store_links(self, links: Iterable[str]) -> int:
        """
        Сохранить/опубликовать набор ссылок.
        Возвращает число обработанных (опубликованных) ссылок.
        """
        raise NotImplementedError


class RedisPublisherAdapter(LinkStorer):
    """
    Адаптер публикации ссылок в Redis Pub/Sub.
    Пароль обязателен при включённом requirepass у Redis, username опционален.
    """
    def __init__(self) -> None:
        self.topic = settings.REDIS_PUBLISH_TOPIC
        self.processed_key = settings.REDIS_PROCESSED_URLS_KEY

        # username передаем только если он задан, чтобы не ломать AUTH
        username = settings.REDIS_USER or None
        password = settings.REDIS_PASSWORD

        self.redis_client = RedisClient(
            host=settings.REDIS_HOST,
            port=settings.REDIS_PORT,
            processed_urls_key=self.processed_key,
            username=username,
            password=password,
        )

    def store_links(self, links: Iterable[str]) -> int:
        """
        Публикует ссылки в канал, полагаясь на логику клиента:
        - проверка processed_urls:zset,
        - публикация только новых,
        - маркировка обработанных.
        Возвращает число реально опубликованных ссылок.
        """
        if not self.redis_client:
            raise ConnectionError("Redis client in adapter is not available.")
        new_links_published = 0
        for item in links:
            url = item.get("url")
            payload = json.dumps({"url": url, "category": item.get("category")}, ensure_ascii=False)
            if self.redis_client.publish(topic=settings.REDIS_PUBLISH_TOPIC, url=url, payload=payload):
                new_links_published += 1
        print(f"Adapter finished. Published {new_links_published} new links to topic '{settings.REDIS_PUBLISH_TOPIC}'.")
        return new_links_published
