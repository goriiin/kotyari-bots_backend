from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    """
    Loads and validates application settings from environment variables.
    """
    # Configure Pydantic to load settings from a .env file.
    model_config = SettingsConfigDict(env_file='./../.env', env_file_encoding='utf-8', extra='ignore')

    #default values
    DZEN_URL_PARSER_HOST: str = "localhost"
    DZEN_URL_PARSER_PORT: int = 8091

    @property
    def construct_server_address(self) -> str:
        return f"{self.DZEN_URL_PARSER_HOST}:{self.DZEN_URL_PARSER_PORT}"

# Create a single, importable instance of the settings to be used across the application.
settings = Settings()