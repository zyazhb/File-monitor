package main

import (
	"flag"
	"io"
	"log"
	"main/notify"

	"github.com/google/logger"
)

//参数名
var (
	h        bool
	f        string
	dir      string
	level    int
	hashflag bool
	server   string
	fan      bool
)

//初始化参数
func init() {
	flag.BoolVar(&h, "h", false, "Show this help")
	flag.StringVar(&f, "f", "", "Choose a file to watch")
	flag.StringVar(&dir, "dir", "", "Choose a dir to monitor")
	flag.IntVar(&level, "level", 1, "Dir level to walk into")
	flag.BoolVar(&hashflag, "hash", false, "Calculate file hash.")
	flag.StringVar(&server, "server", "", "Use rpcreport to Server ip:port")
	flag.BoolVar(&fan, "fan", false, "Use fanotify feature")
}

func main() {
	flag.Parse()
	logger.Init("LoggerExample", true, false, io.Discard)
	logger.SetFlags(log.LstdFlags)
	logger.SetFlags(log.Llongfile)

	if h {
		flag.Usage()
		return
	}

	switch {
	case fan:
		logger.Info("\033[1;30m [*]Fanotify feature is on: " + f + " \033[0m")
		logger.Info("\033[1;30m [*]Watching mountpoint: " + f + " \033[0m")
		mountpoint := ""
		if len(f) != 0 {
			mountpoint = f
		} else if len(dir) != 0 {
			mountpoint = dir
		}
		notify.RunFanotify(mountpoint, hashflag, server)

	case len(f) != 0:
		logger.Info("\033[1;30m [*]Watching file: " + f + " \033[0m")
		filename := []string{f}
		notify.RunInotify(filename, hashflag, server)

	case len(dir) != 0:
		logger.Info("\033[1;30m [*]Start dirwalk: " + dir + " \033[0m")
		notify.RunInotifyForDir(dir, level, hashflag, server)
	}
	flag.Usage()
}
