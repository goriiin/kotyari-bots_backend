import psycopg2
from typing import Dict

class GreenplumWriter:
    """A writer class to handle operations with Greenplum."""

    def __init__(self, db_config: Dict):
        self.tablename = db_config.get("tablename", "parsed_articles")
        self.conn = None
        try:
            print("▶[GreenplumWriter] Попытка подключения к базе данных...")
            self.conn = psycopg2.connect(
                dbname=db_config["dbname"],
                user=db_config["user"],
                password=db_config["password"],
                host=db_config["host"],
                port=db_config["port"],
            )
            self.ensure_table_exists()
            print("[GreenplumWriter] Успешное подключение к Greenplum.")
        except psycopg2.OperationalError as e:
            print(f"[GreenplumWriter] КРИТИЧЕСКАЯ ОШИБКА: Не удалось подключиться к Greenplum: {e}")
            raise

    def ensure_table_exists(self):
        """Ensures the table for parsed articles exists with the correct schema."""
        create_table_query = f"""
        CREATE TABLE IF NOT EXISTS {self.tablename} (
            source_url TEXT NOT NULL PRIMARY KEY,
            title TEXT,
            content TEXT,
            category TEXT,
            parsed_at TIMESTAMPTZ DEFAULT NOW()
        ) DISTRIBUTED BY (source_url);
        """
        try:
            with self.conn.cursor() as cur:
                cur.execute(create_table_query)
                # Ensure other columns/indexes are present if needed
                cur.execute(f"ALTER TABLE {self.tablename} ADD COLUMN IF NOT EXISTS category TEXT;")
                cur.execute(f"CREATE INDEX IF NOT EXISTS idx_{self.tablename}_category ON {self.tablename}(category);")
            self.conn.commit()
        except Exception as e:
            self.conn.rollback()
            print(f"[GreenplumWriter] КРИТИЧЕСКАЯ ОШИБКА при создании/обновлении таблицы: {e}")
            raise

    def insert_article(self, article_data: Dict):
        """Inserts or updates an article in the database."""
        insert_query = f"""
        INSERT INTO {self.tablename} (source_url, title, content, category)
        VALUES (%s, %s, %s, %s)
        ON CONFLICT (source_url) DO UPDATE SET
            title = EXCLUDED.title,
            content = EXCLUDED.content,
            category = COALESCE(EXCLUDED.category, {self.tablename}.category);
        """
        try:
            with self.conn.cursor() as cur:
                cur.execute(
                    insert_query,
                    (
                        article_data.get("source_url"),
                        article_data.get("title"),
                        article_data.get("content"),
                        article_data.get("category"),
                    ),
                )
            self.conn.commit()
        except Exception as e:
            self.conn.rollback()
            print(f"[GreenplumWriter] Ошибка при вставке данных для {article_data.get('source_url')}: {e}")
            raise

    def close(self):
        """Closes the database connection."""
        if self.conn:
            self.conn.close()
            print("[GreenplumWriter] Соединение с Greenplum закрыто.")
