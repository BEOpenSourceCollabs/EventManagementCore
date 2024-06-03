CREATE DATABASE "event-mgmt-db"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

\c event-mgmt-db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";