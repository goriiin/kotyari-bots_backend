#!/bin/bash
set -e

# Запуск основного процесса Greenplum в фоне
/entrypoint.sh &
pid="$!"

echo "Waiting for Greenplum to start..."
# Ожидание доступности порта и ответа от сервера
until psql -h localhost -U gpadmin -d postgres -c '\q' &>/dev/null; do
  sleep 1
done
echo "Greenplum is ready."

# Проверка, выполнялась ли инициализация ранее (идемпотентность)
if [ ! -f /srv/.db_initialized_flag ]; then
  echo "First run: Initializing user and database..."

  # Выполнение SQL для создания роли и БД
  psql -v ON_ERROR_STOP=1 --username "gpadmin" --dbname "postgres" <<-EOSQL
      CREATE ROLE "${GP_USER}" WITH LOGIN PASSWORD '${GP_PASSWORD}';
      CREATE DATABASE "${GP_DB}";
      GRANT ALL PRIVILEGES ON DATABASE "${GP_DB}" TO "${GP_USER}";
EOSQL

  # Создание флага, что инициализация прошла успешно
  touch /srv/.db_initialized_flag
  echo "Initialization complete."
else
  echo "Database already initialized."
fi

# Ожидание завершения основного процесса Greenplum, чтобы контейнер оставался в работе
wait "$pid"
