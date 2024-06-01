# 用户管理系统

完整运行项目的顺序：运行数据库 -> 运行后端 -> 运行前端 -> 在浏览器打开http://localhost:8081/

### 前端运行步骤

在cmd进入fe目录

```sh
cd fe
```

安装依赖

```sh
npm install
```

运行

```sh
npm run serve
```



### 后端运行步骤

在cmd进入be目录

```sh
cd be
```

安装依赖

```sh
go mod tidy
```

运行

```sh
go run main.go
```



### 数据库运行步骤

保证5432端口无程序占用

打开Docker desktop

在cmd进入项目根目录，运行run_db.bat

```sh
./run_db.bat
```



### SQL客户端连接数据库的配置

Host: localhost

Port: 5432

User Name: postgres

Password: postgres

