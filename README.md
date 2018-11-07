### Introduction
-----------------
虽然名字起的比较奇特，但其实就是个爬虫，为了避免资源浪费，并没有使用goroutine并发爬取。
目前只实现了妹子图爬取，默认爬取最近前100条更新，可配置。

特性：
* 爬取妹子图网站，分目录存储
* 支持配置文件
* go语言实现
* 破解防盗链
* 支持http proxy

#### 运行
------

edit config/conf.ini, set proper destination directory, and run commands:

```bash
# prepare config file
cd config
cp conf.ini.example conf.ini
# EDIT conf.ini as you wish

# build
make

# run
./run.sh
```
