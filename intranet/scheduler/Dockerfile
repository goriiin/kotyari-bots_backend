
# 1. Use an official Python runtime as a parent image
FROM python:3.10-slim

# 2. Set the working directory inside the container
WORKDIR /app

ENV PYTHONPATH="/app"

COPY intranet/scheduler/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY api/ /app/api/

COPY intranet/scheduler/ /app/intranet/scheduler/

EXPOSE 8090


CMD ["uvicorn", "intranet.scheduler.main:app", "--host", "0.0.0.0", "--port", "8090"]