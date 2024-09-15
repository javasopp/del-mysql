# 解决mysql启动异常
> 用于解决mysql断电导致的无法启动，本方法很极限，直接删除对应的mysql的数据目录

编译方式: 
```shell
 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go
``` 

windows 编译方式:
```shell
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o main -ldflags="-s -w" main.go
```