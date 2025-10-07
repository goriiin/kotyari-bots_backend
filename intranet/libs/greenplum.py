import psycopg2
from typing import Dict


class GreenplumWriter:
    """
    Класс для управления подключением и записью данных в Greenplum.
    """

    def __init__(self, db_config: Dict):
        self.table_name = db_config.get("table_name", "parsed_articles")
        self.conn = None
        try:
            print("▶[GreenplumWriter] Попытка подключения к базе данных...")
            self.conn = psycopg2.connect(
                dbname=db_config["dbname"],
                user=db_config["user"],
                password=db_config["password"],
                host=db_config["host"],
                port=db_config["port"]
            )
            self._ensure_table_exists()
            print("[GreenplumWriter] Успешное подключение к Greenplum и проверка таблицы.")
        except psycopg2.OperationalError as e:
            print(f"[GreenplumWriter] КРИТИЧЕСКАЯ ОШИБКА: Не удалось подключиться к Greenplum: {e}")
            raise

    def _ensure_table_exists(self):
        """
        Проверяет наличие таблицы для статей и создает ее, если она отсутствует.
        """
        create_table_query = f"""
        CREATE TABLE IF NOT EXISTS {self.table_name} (
            id SERIAL PRIMARY KEY,
            source_url TEXT NOT NULL UNIQUE,
            title TEXT,
            content TEXT,
            parsed_at TIMESTAMPTZ DEFAULT NOW()
        );
        """
        with self.conn.cursor() as cur:
            cur.execute(create_table_query)
            self.conn.commit()

    def insert_article(self, article_data: Dict):
        """
        Вставляет данные статьи в таблицу. Если URL уже существует, ничего не делает.
        """
        insert_query = f"""
        INSERT INTO {self.table_name} (source_url, title, content)
        VALUES (%s, %s, %s)
        ON CONFLICT (source_url) DO NOTHING;
        """
        try:
            with self.conn.cursor() as cur:
                cur.execute(insert_query, (
                    article_data.get("source_url"),
                    article_data.get("title"),
                    article_data.get("content")
                ))
                self.conn.commit()
            print(f"[GreenplumWriter] Данные для URL сохранены: {article_data.get('source_url')}")
        except Exception as e:
            print(f"[GreenplumWriter] Ошибка при вставке данных для URL {article_data.get('source_url')}: {e}")
            self.conn.rollback()  # Откатываем транзакцию в случае ошибки

    def close(self):
        """Закрывает соединение с базой данных."""
        if self.conn:
            self.conn.close()
            print("[GreenplumWriter] Соединение с Greenplum закрыто.")