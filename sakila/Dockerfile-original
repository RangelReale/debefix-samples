# https://daveb.medium.com/using-a-pre-seeded-postgresql-docker-image-for-faster-ci-cd-ccf35305e996
FROM debian:buster-slim AS sakila
# Download the Sakila PostgreSQL example database
RUN apt-get update && apt-get install -y --no-install-recommends wget ca-certificates unzip
RUN wget https://github.com/jOOQ/sakila/archive/refs/heads/main.zip && unzip main.zip 'sakila-main/postgres-sakila-db/*'
# Copy the schema and data .sql files from the Sakila database into
# /docker-entrypoint-initdb.d so they are executed when the
# container initializes
FROM postgres:13 AS sakila-preseed
COPY --from=sakila /sakila-main/postgres-sakila-db/postgres-sakila-schema.sql /docker-entrypoint-initdb.d/01-postgres-sakila-schema.sql
COPY --from=sakila /sakila-main/postgres-sakila-db/postgres-sakila-insert-data.sql /docker-entrypoint-initdb.d/02-postgres-sakila-insert-data.sql
# Modify PostgreSQL's original docker-entrypoint.sh file slightly so
# that PostgreSQL doesn't attempt to start back up after the
#  "initdb" process is complete
RUN grep -v 'exec "$@"' /usr/local/bin/docker-entrypoint.sh > /docker-entrypoint.sh && chmod 755 /docker-entrypoint.sh
# "trust" should never be used in a production setting, but should
# be safe enough for this example
ENV POSTGRES_HOST_AUTH_METHOD trust
ENV POSTGRES_DB sakila
# Run Postgres, it will stop after the .sql files have been installed
# into the database
RUN /docker-entrypoint.sh postgres
# Copy just the pre-seeded PGDATA directory from the previous
# command in our new image
FROM postgres:13
COPY --from=sakila-preseed /var/lib/postgresql/data /var/lib/postgresql/data