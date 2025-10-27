#!/bin/bash
set -e

# Запускаем основной процесс Greenplum в фоне
/entrypoint.sh &
pid="$!"

echo "Waiting for Greenplum to be ready (this may take a moment)..."
# Ожидаем готовности, выполняя psql от имени пользователя gpadmin
until su - gpadmin -c "psql -h localhost -U gpadmin -d postgres -c '\q'" &>/dev/null; do
  sleep 2
done
echo "Greenplum is ready."

# Проверяем, выполнялась ли инициализация ранее
if [ ! -f /srv/.db_initialized_flag ]; then
  echo "First run: Initializing user and database..."

  # Выполняем SQL для создания роли и БД от имени gpadmin
  su - gpadmin -c "psql -v ON_ERROR_STOP=1 --dbname=postgres" <<-EOSQL
      CREATE ROLE "${GP_USER}" WITH LOGIN PASSWORD '${GP_PASSWORD}';
      CREATE DATABASE "${GP_DB}";
      GRANT ALL PRIVILEGES ON DATABASE "${GP_DB}" TO "${GP_USER}";
EOSQL

  # Создаём флаг, что инициализация прошла успешно
  touch /srv/.db_initialized_flag
  echo "Initialization complete."
else
  echo "Database already initialized."
fi

# Ожидаем завершения основного процесса, чтобы контейнер не останавливался
wait "$pid"
