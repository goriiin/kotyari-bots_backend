import time
from abc import ABC, abstractmethod
from typing import List, Generator, Optional

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


class RedisClient(PublisherInterface, SubscriberInterface):
    """
    Реализует оба интерфейса для работы с Redis.
    Инкапсулирует всю логику взаимодействия.
    """

    def __init__(self, host: str, port: int, processed_urls_key: str,
                 username: Optional[str] = None, password: Optional[str] = None):
        self.processed_urls_key = processed_urls_key
        try:
            self._conn = redis.Redis(
                host=host,
                port=port,
                decode_responses=True,
                username=username,
                password=password
            )
            self._conn.ping()
            print(f"[RedisClient] Успешное подключение к Redis ({host}:{port}).")
        except redis.exceptions.AuthenticationError:
            print(f"❌ [RedisClient] КРИТИЧЕСКАЯ ОШИБКА: Неверный логин или пароль для Redis!")
            raise
        except redis.exceptions.ConnectionError as e:
            print(f"❌ [RedisClient] Не удалось подключиться к Redis: {e}")
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


class MockRedisClient(PublisherInterface, SubscriberInterface):
    """
    Мок реализация интерфейсов Publisher и Subscriber для тестирования
    и локальной разработки без реального Redis.
    Использует объекты в памяти для имитации поведения Redis.
    """

    def __init__(self):
        print("✅ [MockRedisClient] Инициализирован фальшивый клиент Redis. Все данные будут храниться в памяти.")
        self._processed_urls = {}  # Имитация ZSET: {url: timestamp}
        self._message_queue = []  # Имитация Pub/Sub очереди
        self._topics = []

    def seed_data(self, initial_urls: List[str]):
        """
        Метод для начального заполнения "очереди" задачами для парсинга.
        :param initial_urls: Список URL для добавления в очередь.
        """
        print(f"[MockRedisClient] Начальное заполнение данными: {len(initial_urls)} URL добавлены в очередь.")
        for url in initial_urls:
            self._message_queue.append({'channel': 'dzen', 'data': url})

    def publish(self, topic: str, url: str) -> bool:
        """Имитирует публикацию, добавляя сообщение в очередь."""
        if url in self._processed_urls:
            print(f"⏭️ [MockPublisher] URL уже существует, публикация отменена: {url}")
            return False

        message = {'channel': topic, 'data': url}
        self._message_queue.append(message)
        print(f"[MockPublisher] URL добавлен в очередь для топика '{topic}': {url}")
        return True

    def listen_for_messages(self, topics: List[str]) -> Generator:
        """
        Имитирует прослушивание. Возвращает сообщения из внутренней очереди,
        а затем "засыпает", имитируя ожидание.
        """
        self._topics = topics
        print(f"🎧 [MockSubscriber] Начал прослушивание топиков: {', '.join(topics)}")

        while True:
            if self._message_queue:
                message = self._message_queue.pop(0)
                if message['channel'] in self._topics:
                    yield message
            else:
                print("...очередь пуста, ожидание...")
                time.sleep(5)

    def mark_as_processed(self, url: str):
        """Имитирует добавление в ZSET, сохраняя URL в словарь."""
        self._processed_urls[url] = time.time()
        print(f"📝 [MockSubscriber] URL помечен как обработанный: {url}")