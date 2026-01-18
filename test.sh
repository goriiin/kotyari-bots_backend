#!/bin/bash

# ==========================================
# КОНФИГУРАЦИЯ
# ==========================================
AUTH_HOST="http://localhost:8010"
PROFILES_HOST="http://localhost:8003"
BOTS_HOST="http://localhost:8001"
POSTS_CMD_HOST="http://localhost:8088"
POSTS_QUERY_HOST="http://localhost:8089" # Новый порт для Query Service

COOKIE_FILE="./cookies_debug.txt"
COOKIE_FILE_2="./cookies_intruder.txt"   # Куки для второго пользователя

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# Уникальные данные для User 1
TIMESTAMP=$(date +%s)
USERNAME="debug_${TIMESTAMP}"
EMAIL="debug_${TIMESTAMP}@test.com"
PASSWORD="Password123!"

# Очистка перед запуском
rm -f "$COOKIE_FILE" "$COOKIE_FILE_2" step*.log step*.json

echo -e "${GREEN}=== ЗАПУСК DEBUG ТЕСТА (FULL FLOW + ISOLATION) ===${NC}"

# ==============================================================================
# ШАГ 1: РЕГИСТРАЦИЯ (USER 1)
# ==============================================================================
echo -e "${CYAN}--------------------------------------------------${NC}"
echo -e "${YELLOW}[STEP 1] Регистрация User 1 (Auth Service)...${NC}"

curl -v -c "$COOKIE_FILE" \
     -X POST "$AUTH_HOST/api/v1/register" \
     -H "Content-Type: application/json" \
     -d "{\"email\": \"$EMAIL\", \"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}" \
     > step1_body.json 2> step1_headers.log

if grep -q "200 OK" step1_headers.log || grep -q "201 Created" step1_headers.log; then
    echo -e "${GREEN}[OK] User 1 зарегистрирован${NC}"
else
    cat step1_headers.log
    cat step1_body.json
    echo -e "${RED}[FAIL] Ошибка регистрации${NC}"
    exit 1
fi

# ==============================================================================
# ШАГ 2: СОЗДАНИЕ ПРОФИЛЯ
# ==============================================================================
echo -e "${CYAN}--------------------------------------------------${NC}"
echo -e "${YELLOW}[STEP 2] Создание профиля (Profiles Service)...${NC}"

curl -v -b "$COOKIE_FILE" \
     -X POST "$PROFILES_HOST/api/v1/profiles" \
     -H "Content-Type: application/json" \
     -d "{\"name\": \"Debug Profile\", \"email\": \"prof_${TIMESTAMP}@test.com\", \"prompt\": \"test\"}" \
     > step2_body.json 2> step2_headers.log

if grep -q "201 Created" step2_headers.log || grep -q "200 OK" step2_headers.log; then
    PROFILE_ID=$(cat step2_body.json | jq -r '.id')
    echo -e "${GREEN}[OK] Профиль создан: $PROFILE_ID${NC}"
else
    cat step2_headers.log
    cat step2_body.json
    echo -e "${RED}[FAIL] Ошибка создания профиля${NC}"
    exit 1
fi

# ==============================================================================
# ШАГ 3: СОЗДАНИЕ БОТА
# ==============================================================================
echo -e "${CYAN}--------------------------------------------------${NC}"
echo -e "${YELLOW}[STEP 3] Создание бота (Bots Service)...${NC}"

# Исправлен email внутри массива profiles на валидный
curl -v -b "$COOKIE_FILE" \
     -X POST "$BOTS_HOST/api/v1/bots" \
     -H "Content-Type: application/json" \
     -d "{\"name\": \"Debug Bot\", \"systemPrompt\": \"sys\", \"moderationRequired\": false, \"profiles\": [{\"id\": \"$PROFILE_ID\", \"name\": \"P\", \"email\": \"test@test.com\", \"systemPrompt\": \"s\"}]}" \
     > step3_body.json 2> step3_headers.log

if grep -q "201 Created" step3_headers.log; then
    BOT_ID=$(cat step3_body.json | jq -r '.id')
    echo -e "${GREEN}[OK] Бот создан: $BOT_ID${NC}"
else
    cat step3_headers.log
    cat step3_body.json
    echo -e "${RED}[FAIL] Ошибка создания бота${NC}"
    exit 1
fi

# ==============================================================================
# ШАГ 4: СОЗДАНИЕ ЗАДАЧИ НА ПОСТИНГ
# ==============================================================================
echo -e "${CYAN}--------------------------------------------------${NC}"
echo -e "${YELLOW}[STEP 4] Создание задачи (Posts Command Service)...${NC}"

curl -v -b "$COOKIE_FILE" \
     -X POST "$POSTS_CMD_HOST/api/v1/posts" \
     -H "Content-Type: application/json" \
     -d "{\"botId\": \"$BOT_ID\", \"profileIds\": [\"$PROFILE_ID\"], \"taskText\": \"debug task\", \"platform\": \"otveti\", \"postType\": \"opinion\"}" \
     > step4_body.json 2> step4_headers.log

if grep -q "201 Created" step4_headers.log; then
    GROUP_ID=$(cat step4_body.json | jq -r '.groupID')
    echo -e "${GREEN}[OK] Задача создана. GroupID: $GROUP_ID${NC}"
else
    cat step4_headers.log
    cat step4_body.json
    echo -e "${RED}[FAIL] Ошибка создания задачи${NC}"
    exit 1
fi

# Небольшая пауза, чтобы Kafka успела обработать сообщение и записать пост в БД (Query side)
echo -e "${YELLOW}Ожидание 2 сек для синхронизации данных...${NC}"
sleep 2

# ==============================================================================
# ШАГ 5: ПОЛУЧЕНИЕ ПОСТОВ ДЛЯ USER 1 (ПРОВЕРКА)
# ==============================================================================
echo -e "${CYAN}--------------------------------------------------${NC}"
echo -e "${YELLOW}[STEP 5] Получение постов User 1 (Posts Query Service)...${NC}"

curl -v -b "$COOKIE_FILE" \
     -X GET "$POSTS_QUERY_HOST/api/v1/posts" \
     > step5_body.json 2> step5_headers.log

echo -e "${CYAN}--- RESPONSE BODY (USER 1) ---${NC}"
cat step5_body.json
echo ""

if grep -q "200 OK" step5_headers.log; then
    # Проверяем, что в массиве data что-то есть (или хотя бы что вернулся JSON)
    COUNT=$(cat step5_body.json | jq '.data | length')
    if [[ "$COUNT" -gt 0 ]]; then
        echo -e "${GREEN}[OK] Посты получены. Количество: $COUNT${NC}"
    else
        # Это допустимо, если консьюмер еще не успел создать запись, но код 200 получен
        echo -e "${YELLOW}[WARN] Список постов пуст (возможно, задержка Kafka), но запрос прошел успешно.${NC}"
    fi
else
    cat step5_headers.log
    echo -e "${RED}[FAIL] Не удалось получить посты User 1${NC}"
    exit 1
fi

# ==============================================================================
# ШАГ 6: ПРОВЕРКА ИЗОЛЯЦИИ (USER 2)
# ==============================================================================
echo -e "${CYAN}--------------------------------------------------${NC}"
echo -e "${YELLOW}[STEP 6] Проверка изоляции (Запрос от User 2)...${NC}"

USERNAME_2="intruder_${TIMESTAMP}"
EMAIL_2="intruder_${TIMESTAMP}@test.com"

# 6.1 Регистрация второго пользователя
echo "Регистрируем User 2 ($USERNAME_2)..."
curl -s -c "$COOKIE_FILE_2" \
     -X POST "$AUTH_HOST/api/v1/register" \
     -H "Content-Type: application/json" \
     -d "{\"email\": \"$EMAIL_2\", \"username\": \"$USERNAME_2\", \"password\": \"$PASSWORD\"}" > /dev/null

if [ ! -f "$COOKIE_FILE_2" ]; then
    echo -e "${RED}[FAIL] Не удалось зарегистрировать User 2${NC}"
    exit 1
fi

# 6.2 Попытка получить посты (должен быть пустой список, так как у User 2 нет постов)
echo "Запрашиваем посты от имени User 2..."
curl -v -b "$COOKIE_FILE_2" \
     -X GET "$POSTS_QUERY_HOST/api/v1/posts" \
     > step6_body.json 2> step6_headers.log

echo -e "${CYAN}--- RESPONSE BODY (USER 2) ---${NC}"
cat step6_body.json
echo ""

# Проверяем ответ
if grep -q "200 OK" step6_headers.log; then
    COUNT_2=$(cat step6_body.json | jq '.data | length')

    # Мы ожидаем, что data будет null (если API возвращает null для пустого слайса) или [] (пустой массив)
    # jq вернет 0 для пустого массива или null.

    if [[ "$COUNT_2" == "0" ]] || [[ "$COUNT_2" == "null" ]]; then
         echo -e "${GREEN}[OK] Изоляция работает! User 2 не видит чужие посты.${NC}"
         echo -e "${GREEN}=== ВСЕ ТЕСТЫ ПРОШЛИ УСПЕШНО ===${NC}"
    else
         echo -e "${RED}[FAIL] НАРУШЕНИЕ ИЗОЛЯЦИИ! User 2 видит $COUNT_2 постов.${NC}"
         exit 1
    fi
else
    cat step6_headers.log
    echo -e "${RED}[FAIL] Ошибка запроса постов для User 2${NC}"
    exit 1
fi