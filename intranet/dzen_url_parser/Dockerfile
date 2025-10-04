
# 1. Use a Debian-based Python image that allows installing system dependencies.
FROM python:3.10-slim-bookworm

# 2. Install system dependencies required for Google Chrome (for Selenium).
# This is the most critical step for the parser to work.
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    # Install Chrome browser
    && wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list \
    && apt-get update && apt-get install -y \
    google-chrome-stable \
    # Clean up APT caches to keep the image size small
    && rm -rf /var/lib/apt/lists/*

# 3. Set the working directory inside the container.
WORKDIR /app

# 4. Set the PYTHONPATH environment variable.
# This is crucial for your monorepo imports (e.g., `from api.protos...`) to work.
ENV PYTHONPATH="/app"

# 5. Copy the shared requirements file first to leverage Docker's layer caching.
# Note the path is relative to the project root.
COPY intranet/scheduler/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 6. Copy the necessary source code directories into the container.
# We need the 'api' directory for the generated proto files.
COPY api/ /app/api/

# Copy this service's specific application code.
COPY intranet/dzen_url_parser/ /app/intranet/dzen_url_parser/

# 7. Expose the gRPC port this service listens on.
EXPOSE 50051

# 8. Define the command to run the gRPC server.
CMD ["python", "-u", "intranet/dzen_url_parser/main.py"]