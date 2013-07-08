##实现


###1.目录结构
```
.
├── conf
│   └── holmes.conf
├── doc
│   ├── image
│   │   ├── deployment.dot
│   │   ├── deployment.png
│   │   ├── holmes.dot
│   │   ├── holmes.png
│   │   ├── syntactical.dot
│   │   ├── syntactical.png
│   │   ├── traffic.dot
│   │   └── traffic.png
│   └── implementation.md
├── format.sh
├── install.sh
├── README.md
└── src
    └── holmes
        ├── accesslog.go
        ├── conf.go
        ├── exporter.go
        ├── filter.go
        ├── main.go
        ├── redis.go
        └── stage.go
```

**accesslog.go**

描述Web服务器访问日志的数据结构以及操作日志的一些函数

**conf.go**

描述holmes配置的数据结构以及操作配置的一些函数

**exporter.go**

过滤结果导出程序，从Redis的队列导出到外部数据目的地（单机文件，HDFS，数据库等等）

**filter.go**

holmes的过滤器

**main.go**

holmes的入口文件

**redis.go**

访问Redis的函数库

**stage**

从数据源（HDFS，单机文件，Nginx的日志输出程序）拉取或接收访问日志，将这些日志存入Redis的日志待处理队列中