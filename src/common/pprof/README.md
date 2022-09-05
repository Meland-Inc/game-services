#cpu 处理器性能: 
	#默认取样30s
		go tool pprof http://ip:port/debug/pprof/profile
	#设置取样45s:
		go tool pprof http://ip:port/debug/pprof/profile\?seconds\=45
#heap 内存性能: 
	#heap 内存性能：
		go tool pprof http://ip:port/debug/pprof/heap
	#设置取样45s:
		go tool pprof http://ip:port/debug/pprof/heap\?second\=45

#常用指令
	topN: 查看耗时或内存占用最多的前N条记录
	web：在浏览器查看火焰图
	web func-name: 查看具体某个函数的性能

default port:6060
go tool pprof http://192.168.1.186:6060/debug/pprof/heap
go tool pprof http://192.168.1.186:6060/debug/pprof/heap?second\=60

go tool pprof http://192.168.1.186:6060//debug/pprof/profile\?seconds\=300