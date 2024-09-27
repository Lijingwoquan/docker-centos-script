# 基于 `docker` 的 `centos` 环境快速搭建


>由于`centos`对用户名的要求 name必须包含英文,不得为纯数字

- 构建并启动centos
```shell
./init.exe new_username new_password
#./init.exe lijingwoquan 123456
```

- 关闭终端后再次进入容器
```shell
docker run -it centos7-yourname:latest
# docker run -it centos7-lijingwoquan:latest  
```

- 若要以管理员的身份运行
```shell
docker run -it -u root centos7-yourname:latest
# docker run -it -u root centos7-lijingwoquan:latest  
```


- 重启容器
```shell
docker restart -it centos7-yourname:latest
```
