FROM postgres:18

# Install pg_cron using the official PostgreSQL repo packages
RUN apt-get update && \
    apt-get install -y postgresql-18-cron && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/lib/postgresql/data && \
    chown -R postgres:postgres /var/lib/postgresql/data