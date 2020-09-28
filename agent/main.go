package main

import (
	"flag"
	"github.com/google/logger"
	"io/ioutil"
	"log"

	"main/fanotify"
)

//参数名
var (
	h        bool
	f        string
	dir      string
	fan      string
	level    int
	daemon   bool
	rpcflag  bool
	hashflag bool
)

//初始化参数
func init() {
	flag.BoolVar(&h, "h", false, "Show this help")
	flag.StringVar(&f, "f", "", "Choose a file to watch")
	flag.StringVar(&dir, "dir", "", "Choose a dir to monitor")
	flag.StringVar(&fan, "fan", "", "Choose a dir to fanotify")
	flag.IntVar(&level, "level", 1, "Dir level to walk into")
	flag.BoolVar(&daemon, "daemon", false, "Start in daemon mode")
	flag.BoolVar(&rpcflag, "rpc", false, "Use rpc report to server")
	flag.BoolVar(&hashflag, "hash", false, "Calculate file hash.")

}

func main() {
	flag.Parse()
	logger.Init("LoggerExample", true, false, ioutil.Discard)
	logger.SetFlags(log.LstdFlags)

	if h {
		flag.Usage()
		return
	}

	switch {
	case len(f) != 0:
		logger.Info("\033[1;30m [*]Watching file: " + f + " \033[0m")
		filename := []string{f}
		inotify(filename, hashflag, rpcflag)

	case len(dir) != 0:
		logger.Info("\033[1;30m [*]Start dirwalk: " + dir + " \033[0m")
		inotifyForDir(dir, level, hashflag, rpcflag)

	case len(fan) != 0:
		logger.Info("\033[1;30m [*]Watching file: " + f + " \033[0m")
		msg, _ := fanotify.Fanotify(fan)
		// if err != nil {
		// 	logger.Error(err)D
		// }
		logger.Info(msg)

	case daemon:
		// rpcreport("aaa")
		return
	}
	flag.Usage()
}
