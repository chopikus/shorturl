-- Database: urldb

-- DROP DATABASE urldb;

CREATE DATABASE urldb
  WITH OWNER = postgres
       ENCODING = 'UTF8'
       TABLESPACE = pg_default
       CONNECTION LIMIT = -1;

-- Table: urls

\c urldb;

-- DROP TABLE urls;

CREATE TABLE urls
(
  url_original VARCHAR(2048) NOT NULL,
  url_code VARCHAR(64) PRIMARY KEY
  CONSTRAINT url_code_constraint CHECK(url_code ~ '^[1-9A-Z]*$')
);


-- Inserting example data into urls table;
-- INSERT INTO urls (url_original, url_code) VALUES ('https://youtube.com', 'ABC11')
