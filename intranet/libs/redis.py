
import time
from abc import ABC, abstractmethod
from typing import List, Generator

import redis

class PublisherInterface(ABC):
    """Интерфейс для публикации новых задач."""

    @abstractmethod
    def publish(self, topic: str, url: str) -> bool:
        """
        Публикует URL в указанный топик, если он еще не был обработан.
        :param topic: Канал (топик) для публикации.
        :param url: URL для проверки и публикации.
        :return: True, если URL был новым и успешно опубликован, иначе False.
        """
        pass


class SubscriberInterface(ABC):
    """Интерфейс для получения и обработки задач."""

    @abstractmethod
    def listen_for_messages(self, topics: List[str]) -> Generator:
        """
        Подписывается на топики и слушает сообщения в режиме long polling.
        Возвращает генератор, который yield'ит сообщения.
        """
        pass

    @abstractmethod
    def mark_as_processed(self, url: str):
        """
        Помечает URL как обработанный (например, после успешного парсинга).
        """
        pass


# --- 2. Конкретная Реализация ---

class RedisClient(PublisherInterface, SubscriberInterface):
    """
    Реализует оба интерфейса для работы с Redis.
    Инкапсулирует всю логику взаимодействия.
    """

    def __init__(self, host: str, port: int, processed_urls_key: str):
        self.processed_urls_key = processed_urls_key
        try:
            self._conn = redis.Redis(host=host, port=port, decode_responses=True)
            self._conn.ping()
            print(f"[RedisClient] Успешное подключение к Redis ({host}:{port}).")
        except redis.exceptions.ConnectionError as e:
            print(f"[RedisClient] Не удалось подключиться к Redis: {e}")
            raise

    def _is_url_processed(self, url: str) -> bool:
        """Внутренний метод для проверки существования URL в ZSET."""
        return self._conn.zscore(self.processed_urls_key, url) is not None

    def publish(self, topic: str, url: str) -> bool:
        """Сначала проверяет URL, затем публикует."""
        if self._is_url_processed(url):
            print(f"[Publisher] URL уже существует в ZSET, публикация отменена: {url}")
            return False

        subscribers_count = self._conn.publish(topic, url)
        print(f"[Publisher] URL опубликован в топик '{topic}'. Слушателей: {subscribers_count}. URL: {url}")
        return True

    def listen_for_messages(self, topics: List[str]) -> Generator:
        """Подписывается на топики и возвращает генератор сообщений."""
        pubsub = self._conn.pubsub()
        if not topics:
            print("[Subscriber] Нет топиков для подписки.")
            return

        pubsub.subscribe(*topics)
        print(f"[Subscriber] Подписан на топики: {', '.join(topics)}")

        for message in pubsub.listen():
            if message['type'] == 'message':
                yield message

    def mark_as_processed(self, url: str):
        """Помечает URL как обработанный в ZSET."""
        self._conn.zadd(self.processed_urls_key, {url: time.time()})
        print(f"[Subscriber] URL помечен как обработанный: {url}")