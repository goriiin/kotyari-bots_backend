from pydantic_settings import BaseSettings, SettingsConfigDict
from typing import List

class Settings(BaseSettings):
    """
    Loads and validates application settings from environment variables.
    """
    # Configure Pydantic to load settings from a .env file.
    model_config = SettingsConfigDict(env_file='./../../.env', env_file_encoding='utf-8', extra='ignore')

    #default values
    DZEN_URL_PARSER_HOST: str = "localhost"
    DZEN_URL_PARSER_PORT: int = 50051

    REDIS_HOST: str = "localhost"
    REDIS_PORT: int = 6379


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

    DZEN_SCROLL_COUNT: int = 5
    DZEN_SCROLL_DELAY_SECONDS: int = 3

    @property
    def construct_server_address(self) -> str:
        return f"{self.DZEN_URL_PARSER_HOST}:{self.DZEN_URL_PARSER_PORT}"

    @property
    def construct_redis_address(self) -> str:
        return f"{self.REDIS_HOST}:{self.REDIS_PORT}"

# Create a single, importable instance of the settings to be used across the application.
settings = Settings()