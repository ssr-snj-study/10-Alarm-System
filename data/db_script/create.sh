#!/bin/bash
set -e

PGPASSWORD=$POSTGRESQL_PASSWORD psql -v ON_ERROR_STOP=1 --username "$POSTGRESQL_USERNAME" --dbname "$POSTGRESQL_DATABASE" <<-EOSQL
CREATE TABLE user_tb (
    user_id serial4 NOT NULL ,
    email varchar ,
    country_code int4 ,
    phone_number int4 ,
    created_art timestamp NOT NULL ,
  CONSTRAINT user_pkey PRIMARY KEY (user_id)
  );

  CREATE TABLE device (
    id serial4 NOT NULL,
    device_token varchar,
    user_id int4,
    last_logged_in_at timestamp NOT NULL,
  CONSTRAINT device_pkey PRIMARY KEY (id)
  );

  CREATE TABLE message (
      id serial4 NOT NULL,
      user_id int4,
      send_time timestamp NOT NULL,
      contents text NOT NULL,
      receiver varchar NOT NULL,
    CONSTRAINT message_pkey PRIMARY KEY (id)
  );
EOSQL