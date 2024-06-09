CREATE DATABASE admin_system_db;
CREATE USER admin_system_user WITH PASSWORD '766515';
GRANT ALL PRIVILEGES ON DATABASE admin_system_db TO admin_system_user;
\c admin_system_db
CREATE SCHEMA admin_system;
GRANT CREATE ON SCHEMA admin_system TO admin_system_user;
ALTER SCHEMA admin_system OWNER TO admin_system_user;
ALTER ROLE admin_system_user SET search_path TO admin_system, public;
exit
psql -U admin_system_user -d admin_system_db
CREATE TABLE IF NOT EXISTS admin_system.admin (id SERIAL,name VARCHAR(255),email VARCHAR(255),password VARCHAR(255),phone_number VARCHAR(255),avatar_filename VARCHAR(255),table_name VARCHAR(255),remarks VARCHAR(255));