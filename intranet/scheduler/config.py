from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    """
    Loads and validates application settings from environment variables.
    """
    model_config = SettingsConfigDict(env_file='./../.env', env_file_encoding='utf-8', extra='ignore')

    # defaults tuned for docker network
    DZEN_URL_PARSER_HOST: str = "dzen-parser"
    DZEN_URL_PARSER_PORT: int = 8091

    # increase time budget for long Selenium session
    GRPC_TIMEOUT_SECONDS: int = 480
    GRPC_RETRY_MAX: int = 3
    GRPC_RETRY_BACKOFF_MS: int = 500

    GRPC_KEEPALIVE_TIME_MS: int = 45000
    GRPC_KEEPALIVE_TIMEOUT_MS: int = 20000
    GRPC_KEEPALIVE_PERMIT_WITHOUT_CALLS: int = 0
    GRPC_MAX_PINGS_WITHOUT_DATA: int = 0

    @property
    def construct_server_address(self) -> str:
        return f"{self.DZEN_URL_PARSER_HOST}:{self.DZEN_URL_PARSER_PORT}"

settings = Settings()
