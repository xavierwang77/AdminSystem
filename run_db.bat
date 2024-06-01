@echo off
REM 这个脚本用于运行Postgres数据库容器

REM 定义变量
set POSTGRES_CONTAINER_NAME=admin_system_pg
set POSTGRES_USER=postgres
set POSTGRES_PASSWORD=postgres
set POSTGRES_PORT=5432
set ADMIN_SYSTEM_DB=admin_system_db
set ADMIN_SYSTEM_USER=admin_system_user
set ADMIN_SYSTEM_USER_PASSWORD=766515
set ADMIN_SYSTEM_SCHEMA=admin_system

REM 检查Docker是否已经安装
where docker >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Docker is not installed. Please install Docker first.
    exit /b
)

REM 检查是否存在postgres镜像
docker images -q postgres:latest >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Postgres image not found. Pulling the image...
    docker pull postgres:latest
)

REM 检查是否已有同名的容器
docker ps -a | findstr /i %POSTGRES_CONTAINER_NAME% >nul
if %ERRORLEVEL% equ 0 (
    echo Postgres container already exists. Starting the container...
    docker start %POSTGRES_CONTAINER_NAME%
) else (
    echo Creating a new Postgres container...
    docker run --name %POSTGRES_CONTAINER_NAME% -e POSTGRES_PASSWORD=%POSTGRES_PASSWORD% -p %POSTGRES_PORT%:5432 -d postgres:latest
)

timeout /t 5

REM 输出容器运行状态
if %ERRORLEVEL% equ 0 (
    echo Postgres container is running.
    echo Container name: %POSTGRES_CONTAINER_NAME%
    echo Username: %POSTGRES_USER%
    echo Password: %POSTGRES_PASSWORD%
    echo Port: %POSTGRES_PORT%
) else (
    echo Failed to start Postgres container.
    exit /b 1
)

REM 连接到Postgres数据库并创建数据库和用户
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -c "CREATE DATABASE %ADMIN_SYSTEM_DB%;"
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -c "CREATE USER %ADMIN_SYSTEM_USER% WITH PASSWORD '%ADMIN_SYSTEM_USER_PASSWORD%';"
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -c "GRANT ALL PRIVILEGES ON DATABASE %ADMIN_SYSTEM_DB% TO %ADMIN_SYSTEM_USER%;"

REM 创建模式并授权给用户
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -d %ADMIN_SYSTEM_DB% -c "CREATE SCHEMA %ADMIN_SYSTEM_SCHEMA%;"
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -d %ADMIN_SYSTEM_DB% -c "GRANT CREATE ON SCHEMA %ADMIN_SYSTEM_SCHEMA% TO %ADMIN_SYSTEM_USER%;"
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -d %ADMIN_SYSTEM_DB% -c "ALTER SCHEMA %ADMIN_SYSTEM_SCHEMA% OWNER TO %ADMIN_SYSTEM_USER%;"
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %POSTGRES_USER% -h localhost -p %POSTGRES_PORT% -d %ADMIN_SYSTEM_DB% -c "ALTER ROLE %ADMIN_SYSTEM_USER% SET search_path TO %ADMIN_SYSTEM_SCHEMA%, public;"

REM 执行SQL文件
docker exec -it %POSTGRES_CONTAINER_NAME% psql -U %ADMIN_SYSTEM_USER% -h localhost -p %POSTGRES_PORT% -d %ADMIN_SYSTEM_DB% -c "CREATE TABLE IF NOT EXISTS admin_system.admin (id SERIAL,name VARCHAR(255),email VARCHAR(255),password VARCHAR(255),phone_number VARCHAR(255),avatar_filename VARCHAR(255),table_name VARCHAR(255),remarks VARCHAR(255));"