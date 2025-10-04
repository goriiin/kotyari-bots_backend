from abc import ABC, abstractmethod
from typing import List, Dict, Union

class AbstractBaseParser(ABC):
    """
    Абстрактный базовый класс для ВСЕХ парсеров.
    Определяет единый контракт: каждый парсер должен иметь метод `parse`.
    """
    @abstractmethod
    def parse(self, target: str) -> Union[Dict, List[Dict]]:
        """
        Основной метод парсинга.
        Принимает 'target' (это может быть URL, ID группы, поисковый запрос и т.д.).
        Возвращает либо один словарь с данными (для парсинга одной страницы),
        либо список словарей (для API, возвращающих много постов).
        """
        pass

from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from webdriver_manager.chrome import ChromeDriverManager

class BaseParser(ABC):
    @abstractmethod
    def parse(self, target: str) -> Union[Dict, List[Dict]]:
        pass

class BaseBrowserParser(BaseParser, ABC):
    def __init__(self):
        print("▶[BaseBrowserParser] Инициализация... Попытка создать WebDriver.")
        self.driver = self._get_webdriver()
        print("[BaseBrowserParser] WebDriver успешно создан.")

    def _get_webdriver(self):
        options = webdriver.ChromeOptions()
        options.add_argument("--headless")
        options.add_argument("--no-sandbox")
        options.add_argument("--disable-dev-shm-usage")
        options.add_argument("--log-level=3")
        options.add_argument(
            "user-agent=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
        options.add_experimental_option("prefs", {"profile.managed_default_content_settings.images": 2})

        try:
            print("   [WebDriver] Попытка использовать ChromeDriverManager для локального запуска.")
            driver = webdriver.Chrome(service=Service(ChromeDriverManager().install()), options=options)
            print("   [WebDriver] Успешно. Установлен и используется драйвер через ChromeDriverManager.")
            return driver
        except Exception as e:

            print(f"   [WebDriver] КРИТИЧЕСКАЯ ОШИБКА: Не удалось создать драйвер: {e}")
            raise

    def close(self):
        if hasattr(self, 'driver') and self.driver:
            self.driver.quit()


import requests

class BaseApiParser(BaseParser, ABC):
    """
    Базовый класс для парсеров, которые обращаются к API.
    Может содержать общую логику для HTTP-запросов, например, сессию.
    """
    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({'Accept': 'application/json'})