FROM pgduckdb/pgduckdb:18-main

# Switch to root-user to allow package installation. By default the image is 
# set to non-root user
USER root

# pg_cron setup
RUN apt-get update && apt-get install -y \
    postgresql-18-cron \
    && rm -rf /var/lib/apt/lists/*

# Update the existing pg_duckdb configuration to include pg_cron
RUN sed -i "s/shared_preload_libraries='pg_duckdb'/shared_preload_libraries='pg_duckdb, pg_cron'/" /usr/share/postgresql/postgresql.conf.sample \
    && echo "cron.database_name = 'postgres'" >> /usr/share/postgresql/postgresql.conf.sample

# Switch back to non-root user
# USER postgres
