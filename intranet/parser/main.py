# /intranet/parser/main.py
from intranet.libs.redis import MockRedisClient, RedisClient
from intranet.parser.config import config
from intranet.parser.parsers.dzen import DzenParser

INITIAL_URLS_FOR_MOCK = [
    "https://dzen.ru/a/aIIhDlT2Qh9H9FOZ",
    "https://dzen.ru/a/aN4o2LNMgXBpOxFH"
]

PARSER_MAPPING = {
    "dzen": DzenParser,
}


class ParserOrchestrator:
    """
    Главный класс-оркестратор. Выполняет задачи по парсингу последовательно.
    """

    def __init__(self):
        if config.get('use_mock_redis'):
            mock_client = MockRedisClient()
            mock_client.seed_data(INITIAL_URLS_FOR_MOCK)
            self.redis_client = mock_client
        else:
            self.redis_client = RedisClient(
                host=config['redis']['host'],
                port=config['redis']['port'],
                processed_urls_key=config['redis'].get('processed_urls_key', 'processed_urls:zset'),
                username=config['redis'].get('user'),
                password=config['redis'].get('password')
            )
        print("Оркестратор запущен в однопоточном (синхронном) режиме.")

    def process_item(self, item_data: dict):
        """Обрабатывает один элемент (пост/статью) после успешного парсинга."""
        url = item_data.get("source_url")
        if not url:
            print("⚠️ Получены данные без source_url, невозможно отметить как обработанные.")
            return

        self.redis_client.mark_as_processed(url)

        print("\n" + "=" * 70)
        print(f"УСПЕШНО СПАРСЕНА СТАТЬЯ")
        print(f"   ЗАГОЛОВОК: {item_data.get('title')}")
        print(f"   URL: {url}")
        print(f"   КОНТЕНТ (первые 150 символов): {item_data.get('content', '')[:150].replace(chr(10), ' ')}...")
        print("=" * 70 + "\n")

    def run(self):
        """
        Главный цикл приложения. Последовательно получает и обрабатывает задачи.
        Блокируется на время выполнения каждой задачи.
        """
        topics = config.get('active_topics', [])

        print(f"Начинаю прослушивание топиков: {', '.join(topics) if topics else 'Нет топиков для прослушивания'}")
        for message in self.redis_client.listen_for_messages(topics):
            topic = message['channel']
            target = message['data']
            parser_class = PARSER_MAPPING.get(topic)

            if not parser_class:
                print(f"Для топика '{topic}' не найден соответствующий парсер. Задача пропущена.")
                continue

            print(f"🔹 Взята в работу задача: {target}")
            parser_instance = None
            try:
                parser_instance = parser_class()
                results = parser_instance.parse(target)

                if isinstance(results, list):
                    for item in results:
                        if item.get("status") == "success":
                            self.process_item(item)
                elif isinstance(results, dict):
                    if results.get("status") == "success":
                        self.process_item(results)

            except Exception as e:
                print(f"Критическая ошибка при обработке {target}: {e}")
            finally:
                if hasattr(parser_instance, 'close'):
                    parser_instance.close()
                print(f"Задача завершена: {target}. Ожидание следующей...")


if __name__ == "__main__":
    configured_topics = set(config.get('active_topics', []))
    registered_parsers = set(PARSER_MAPPING.keys())

    if not configured_topics.issubset(registered_parsers):
        missing = configured_topics - registered_parsers
        print(f"КРИТИЧЕСКАЯ ОШИБКА: В config.yml указаны топики, для которых нет парсеров: {missing}")
        exit(1)

    orchestrator = ParserOrchestrator()
    orchestrator.run()