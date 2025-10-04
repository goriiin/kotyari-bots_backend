import os
import yaml
from dotenv import load_dotenv


def load_config():
    """
    Загружает конфигурацию из нескольких источников и объединяет их.
    Порядок приоритета (поздние перезаписывают ранние):
    1. config.yml (базовые настройки, которые можно коммитить в Git).
    2. .env файл (секреты и локальные переопределения).
    3. Системные переменные окружения.
    """
    load_dotenv()

    try:
        with open('../config.yml', 'r') as f:
            config = yaml.safe_load(f)
    except FileNotFoundError:
        print("❌ КРИТИЧЕСКАЯ ОШИБКА: Файл config.yml не найден в корне проекта `intranet/`!")
        exit(1)
    except yaml.YAMLError as e:
        print(f"❌ КРИТИЧЕСКАЯ ОШИБКА: Ошибка парсинга config.yml: {e}")
        exit(1)

    # redis / mock
    config['use_mock_redis'] = os.getenv('USE_MOCK_REDIS', 'false').lower() in ('true', '1', 't')

    if 'redis' not in config:
        config['redis'] = {}
    config['redis']['host'] = os.getenv("REDIS_HOST", config['redis'].get('host', 'localhost'))
    config['redis']['port'] = int(os.getenv("REDIS_PORT", config['redis'].get('port', 6379)))

    config['redis']['user'] = os.getenv("REDIS_USER")
    config['redis']['password'] = os.getenv("REDIS_PASSWORD")

    return config


config = load_config()