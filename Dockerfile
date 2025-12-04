FROM postgres:18

# 1. Update package lists
# 2. Install the pg_cron extension for this specific Postgres version
# 3. Clean up apt cache to keep image small
RUN apt-get update \
    && apt-get install -y postgresql-18-cron \
    && rm -rf /var/lib/apt/lists/*
