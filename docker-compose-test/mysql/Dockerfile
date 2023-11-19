#使用mysql8.0在后续的sql自动执行中会有问题，所以这里采用mysql5.7
FROM mysql:5.7.25

MAINTAINER my-mysql


#复制初始化sql到容器内部
COPY ./test.sql /usr/local/test.sql

#把要执行的shell文件放到/docker-entrypoint-initdb.d/目录下，容器会自动执行这个shell
COPY ./docker-entrypoint.sh /docker-entrypoint-initdb.d/docker-entrypoint.sh

#给执行文件增加可执行权限
RUN chmod a+x /docker-entrypoint-initdb.d/docker-entrypoint.sh

#把sql
ADD test.sql /docker-entrypoint-initdb.d

#设置容器启动时执行的命令
CMD sh /docker-entrypoint-initdb.d/docker-entrypoint.sh