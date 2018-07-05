# gfs
一个由Go语言实现的分布式文件系统

gfs/common包含了其他模块需要使用的通用的东西
gfs/gfssdk 是一个sdk，可以用于与gfs通讯
gfs/gfsclient 是一个客户端
gfs/gfsmaster 是gfs的master
gfs/gfsnode 是gfs的datanode

系统有如下约定：所有的http请求都使用Post，参数放置在请求体中，URL中不包含参数；


gfs/common中包含session与cookie的实现、文件系统FileSystem的抽象、
一些工具、GFS的配置与读取。