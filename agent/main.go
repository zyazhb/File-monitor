package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/google/logger"
)

//参数名
var (
	h        bool
	f        string
	dir      string
	level    int
	daemon   bool
	rpcflag  bool
	hashflag bool
	server   string
)

//初始化参数
func init() {
	flag.BoolVar(&h, "h", false, "Show this help")
	flag.StringVar(&f, "f", "", "Choose a file to watch")
	flag.StringVar(&dir, "dir", "", "Choose a dir to monitor")
	flag.IntVar(&level, "level", 1, "Dir level to walk into")
	flag.BoolVar(&daemon, "daemon", false, "Start in daemon mode")
	flag.BoolVar(&rpcflag, "rpc", false, "Use rpc report to server")
	flag.BoolVar(&hashflag, "hash", false, "Calculate file hash.")
	flag.StringVar(&server, "server", "", "Server ip:port")
}

func main() {
	flag.Parse()
	logger.Init("LoggerExample", true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)
	logger.SetFlags(log.Llongfile)

	if h {
		flag.Usage()
		return
	}

	switch {
	case len(f) != 0:
		logger.Info("\033[1;30m [*]Watching file: " + f + " \033[0m")
		filename := []string{f}
		inotify(filename, hashflag, rpcflag, server)

	case len(dir) != 0:
		logger.Info("\033[1;30m [*]Start dirwalk: " + dir + " \033[0m")
		inotifyForDir(dir, level, hashflag, rpcflag, server)
	case daemon:
		// rpcreport("aaa")
		return
	}
	flag.Usage()
}
