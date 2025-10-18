from typing import List, Optional
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    # .env рядом с compose в intranet/, код читает переменные без переименований
    model_config = SettingsConfigDict(
        env_file='./../.env',
        env_file_encoding='utf-8',
        extra='ignore',
    )

    # gRPC сервера
    DZEN_URL_PARSER_HOST: str = "dzen-parser"
    DZEN_URL_PARSER_PORT: int = 8091

    # Redis
    REDIS_HOST: str = "redis"
    REDIS_PORT: int = 6379
    # Делаем имя пользователя необязательным, чтобы не падать при отсутствующем ACL
    REDIS_USER: Optional[str] = None
    REDIS_PASSWORD: Optional[str] = None
    REDIS_PROCESSED_URLS_KEY: str = "processed_urls:zset"
    REDIS_PUBLISH_TOPIC: str = "dzen"

    DZEN_ARTICLES_URLs: List[str]  = [
        "https://dzen.ru/topic/travel",
        "https://dzen.ru/articles",
        "https://dzen.ru/topic/food",
        "https://dzen.ru/topic/culture",
        "https://dzen.ru/topic/economy",
        "https://dzen.ru/topic/it",
        "https://dzen.ru/topic/auto",
        "https://dzen.ru/topic/games",
        "https://dzen.ru/topic/pets",
        "https://dzen.ru/topic/multiki",
        "https://dzen.ru/topic/science"
    ]

    # Скроллинг Selenium
    DZEN_SCROLL_COUNT: int = 3
    DZEN_SCROLL_DELAY_SECONDS: int = 3

    # Тайминги Selenium/WebDriver ожиданий (сек)
    SELENIUM_PAGELOAD_TIMEOUT: int = 30
    SELENIUM_IMPLICIT_WAIT: int = 5

    # Диагностика
    DEBUG_LOG_SELECTORS: bool = True

    @property
    def construct_server_address(self) -> str:
        return f"{self.DZEN_URL_PARSER_HOST}:{self.DZEN_URL_PARSER_PORT}"


settings = Settings()
