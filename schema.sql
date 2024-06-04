CREATE DATABASE "event-mgmt-db"
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

\c event-mgmt-db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role AS ENUM ('user', 'admin', 'organizer');

CREATE TABLE IF NOT EXISTS users (
  id uuid DEFAULT gen_random_uuid(),
  username VARCHAR(20) NOT NULL,
  email VARCHAR(256) NOT NULL,
  password TEXT NOT NULL,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  birth_date DATE,
  role role,
  verified boolean,
  about VARCHAR(500),
  PRIMARY KEY (id)
);