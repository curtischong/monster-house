# Postgres Client
This client is used to save and fetch data from the Postgres client running in docker

# Notes:
- Since Redshift support in localstack is not supported, we are using postgres instead.
  See https://github.com/localstack/localstack/issues/833

# TODO
- Run migrations after docker compose after initialization script:
  Initialization scripts