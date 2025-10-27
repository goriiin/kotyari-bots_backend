from abc import ABC, abstractmethod
from typing import List, Dict, Union, Optional

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

from intranet.libs.driver import create_anti_detect_driver

class BaseParser(ABC):
    @abstractmethod
    def parse(self, target: str) -> Union[Dict, List[Dict]]:
        pass

class BaseBrowserParser(BaseParser, ABC):
    def __init__(self, proxy: Optional[str] = None):
        print("▶[BaseBrowserParser] Инициализация... Попытка создать anti-detect WebDriver.")
        self.driver = create_anti_detect_driver(proxy=proxy)
        print("[BaseBrowserParser] WebDriver успешно создан.")

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
