from typing import Dict
from selenium.common.exceptions import NoSuchElementException, TimeoutException
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from .interface import BaseBrowserParser


class DzenParser(BaseBrowserParser):
    WAIT_TIMEOUT = 15

    def parse(self, url: str) -> Dict:
        print(f"⚙️ [DzenParser] Начинаю парсинг: {url}")
        try:
            self.driver.get(url)
            wait = WebDriverWait(self.driver, self.WAIT_TIMEOUT)

            title_element = wait.until(
                EC.visibility_of_element_located((By.XPATH, "//h1[@data-testid='article-title']"))
            )
            title = title_element.text
            # Ждем появления контейнера с текстом статьи
            article_body = wait.until(
                EC.visibility_of_element_located((By.CSS_SELECTOR, "div[itemprop='articleBody']"))
            )

            # Ищем все параграфы внутри найденного контейнера
            p_elements = article_body.find_elements(By.CSS_SELECTOR, "p[data-article-block='true']")
            article_text = "\n".join([p.text for p in p_elements if p.text.strip()])
            return {
                "source_url": url,
                "title": title,
                "content": article_text,
                "status": "success"
            }

        except TimeoutException:
            print(f"❌ [DzenParser] Таймаут ожидания элементов на странице (возможно, CAPTCHA): {url}")
            return {"source_url": url, "error": f"TimeoutException after {self.WAIT_TIMEOUT}s", "status": "failed"}

        except NoSuchElementException:
            print(f"⚠️ [DzenParser] Не найдены ключевые элементы на странице: {url}")
            return {"source_url": url, "error": "Content not found", "status": "failed"}

        except Exception as e:
            print(f"❌ [DzenParser] Неизвестная ошибка при обработке {url}: {e}")
            return {"source_url": url, "error": str(e), "status": "failed"}

        finally:
            self.close()