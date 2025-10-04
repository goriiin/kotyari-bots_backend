import threading
from concurrent.futures import ThreadPoolExecutor

from intranet.libs.redis import MockRedisClient, RedisClient
from intranet.parser.config import config
from intranet.parser.parsers.dzen import DzenParser
from intranet.libs.greenplum import GreenplumWriter

INITIAL_URLS_FOR_MOCK = [
    "https://dzen.ru/a/aIIhDlT2Qh9H9FOZ",
    "https://dzen.ru/a/aN4o2LNMgXBpOxFH"
]
PARSER_MAPPING = {"dzen": DzenParser}


class ParserOrchestrator:
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

        self.storage_writer = GreenplumWriter(config['greenplum'])
        self.executor = ThreadPoolExecutor(max_workers=config['parser']['max_workers'])
        print(f"Оркестратор запущен с {config['parser']['max_workers']} воркерами.")

    def process_item(self, item_data: dict):
        url = item_data.get("source_url")
        if not url:
            print("Получены данные без source_url.")
            return

        self.redis_client.mark_as_processed(url)
        self.storage_writer.insert_article(item_data)
        print(f"УСПЕШНО ОБРАБОТАНА И СОХРАНЕНА СТАТЬЯ: {item_data.get('title')}")

    def worker(self, parser_class, target: str):
        thread_id = threading.get_ident()
        print(f"Воркер [{thread_id}] взял в работу таргет: {target}")
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
        except Exception as err:
            print(f"Критическая ошибка при обработке {target}: {err}")
        finally:
            if hasattr(parser_instance, 'close'):
                parser_instance.close()

    def run(self):
        topics = config.get('active_topics', [])
        print(f"Начинаю прослушивание топиков: {', '.join(topics) if topics else 'Нет топиков для прослушивания'}")
        for message in self.redis_client.listen_for_messages(topics):
            topic = message['channel']
            target = message['data']
            parser_class = PARSER_MAPPING.get(topic)
            if parser_class:
                self.executor.submit(self.worker, parser_class, target)
            else:
                print(f"Для топика '{topic}' не найден соответствующий парсер.")


if __name__ == "__main__":
    configured_topics = set(config.get('active_topics', []))
    registered_parsers = set(PARSER_MAPPING.keys())
    if not configured_topics.issubset(registered_parsers):
        missing = configured_topics - registered_parsers
        print(f"КРИТИЧЕСКАЯ ОШИБКА: В config.yml указаны топики, для которых нет парсеров: {missing}")
        exit(1)

    try:
        orchestrator = ParserOrchestrator()
        orchestrator.run()
    except Exception as e:
        print(f"КРИТИЧЕСКАЯ ОШИБКА НА СТАРТЕ ПРИЛОЖЕНИЯ: {e}")
        exit(1)