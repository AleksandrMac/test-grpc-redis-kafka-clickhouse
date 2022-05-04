-- CREATE USER testhezzl_user PASSWORD 'testpass';
CREATE DATABASE testhezzl_db
    WITH
    OWNER testhezzl_user
    ENCODING = 'UTF8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = 50;