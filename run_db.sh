#!/bin/bash
# 这个脚本用于运行Postgres数据库容器

# 定义变量
POSTGRES_CONTAINER_NAME="admin_system_pg"
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="postgres"
POSTGRES_PORT="5432"
ADMIN_SYSTEM_DB="admin_system_db"
ADMIN_SYSTEM_USER="admin_system_user"
ADMIN_SYSTEM_USER_PASSWORD="766515"
ADMIN_SYSTEM_SCHEMA="admin_system"

# 检查Docker是否已经安装
if ! command -v docker &> /dev/null
then
    echo "Docker is not installed. Please install Docker first."
    exit
fi

# 检查是否存在postgres镜像
if [[ "$(docker images -q postgres:latest 2> /dev/null)" == "" ]]; then
    echo "Postgres image not found. Pulling the image..."
    docker pull postgres:latest
fi

# 检查是否已有同名的容器
if [[ "$(docker ps -a | grep $POSTGRES_CONTAINER_NAME)" != "" ]]; then
    echo "Postgres container already exists. Starting the container..."
    docker start $POSTGRES_CONTAINER_NAME
else
    echo "Creating a new Postgres container..."
    docker run --name $POSTGRES_CONTAINER_NAME -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD -p $POSTGRES_PORT:5432 -d postgres:latest
fi

# 输出容器运行状态
if [ $? -eq 0 ]; then
  echo "Postgres container is running."
  echo "Container name: $POSTGRES_CONTAINER_NAME"
  echo "Database name: $ADMIN_SYSTEM_DB"
  echo "Username: $POSTGRES_USER"
  echo "Password: $POSTGRES_PASSWORD"
  echo "Port: $POSTGRES_PORT"
else
  echo "Failed to start Postgres container."
  exit 1
fi

# 连接到Postgres数据库并创建数据库和用户
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -c "CREATE DATABASE $ADMIN_SYSTEM_DB;"
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -c "CREATE USER $ADMIN_SYSTEM_USER WITH PASSWORD '$ADMIN_SYSTEM_USER_PASSWORD';"
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -c "GRANT ALL PRIVILEGES ON DATABASE $ADMIN_SYSTEM_DB TO $ADMIN_SYSTEM_USER;"

# 创建模式并授权给用户
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $ADMIN_SYSTEM_DB -c "CREATE SCHEMA $ADMIN_SYSTEM_SCHEMA;"
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $ADMIN_SYSTEM_DB -c "GRANT CREATE ON SCHEMA $ADMIN_SYSTEM_SCHEMA TO $ADMIN_SYSTEM_USER;"
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $ADMIN_SYSTEM_DB -c "ALTER SCHEMA $ADMIN_SYSTEM_SCHEMA OWNER TO $ADMIN_SYSTEM_USER;"
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $POSTGRES_USER -d $ADMIN_SYSTEM_DB -c "ALTER ROLE $ADMIN_SYSTEM_USER SET search_path TO $ADMIN_SYSTEM_SCHEMA, public;"

# 执行SQL文件
docker exec -it $POSTGRES_CONTAINER_NAME psql -U $ADMIN_SYSTEM_USER -d $ADMIN_SYSTEM_DB -f admin_table.sql