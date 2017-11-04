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

```bash
cd tissueabsuer/test 
go run test.go

>
temporarily change working dir to: /Users/bruce/Code/go/src/tissueabsuer/test/../config
change working dir back: /Users/bruce/Code/go/src/tissueabsuer/test
listening for albums...
start crawling mzitu: http://mzitu.com/all
got album: 沙滩尤物艺轩 海风撩起的不只她的长发,还有春心荡漾 http://www.mzitu.com/107734
trying to build album dir
dir created successfully
got image link: http://i.meizitu.net/2017/11/02a01.jpg
storing file: /Users/bruce/Downloads/Image/沙滩尤物艺轩 海风撩起的不只她的长发,还有春心荡漾/02a01.jpg
got image link: http://i.meizitu.net/2017/11/02a02.jpg
storing file: /Users/bruce/Downloads/Image/沙滩尤物艺轩 海风撩起的不只她的长发,还有春心荡漾/02a02.jpg
got image link: http://i.meizitu.net/2017/11/02a03.jpg
storing file: /Users/bruce/Downloads/Image/沙滩尤物艺轩 海风撩起的不只她的长发,还有春心荡漾/02a03.jpg
got image link: http://i.meizitu.net/2017/11/02a04.jpg
storing file: /Users/bruce/Downloads/Image/沙滩尤物艺轩 海风撩起的不只她的长发,还有春心荡漾/02a04.jpg
....

```
