-- drop\create user
DROP USER IF EXISTS snippetbox;
CREATE USER snippetbox
    WITH PASSWORD 'P@ssw0rd';

-- drop\create db
DROP DATABASE IF EXISTS snippetbox;
CREATE DATABASE snippetbox
    WITH OWNER snippetbox;
