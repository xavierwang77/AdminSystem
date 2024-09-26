# 使用官方的 Nginx 镜像作为基础镜像
FROM nginx:latest

# 将本地的 HTML 文件复制到编译目录
COPY ./dist /var/deploy

# 将本地nginx配置文件复制到对应目录
COPY ./admin_fe.conf /etc/nginx/conf.d

# 暴露端口 6230
EXPOSE 6230

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]