package main

import (
	"flag"
	"log"
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

}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	switch {
	case len(f) != 0:
		log.Print("[+]Watching file: " + f)
		filename := []string{f}
		inotify(filename, hashflag, rpcflag)

	case len(dir) != 0:
		log.Print("[+]Start dirwalk: " + dir)
		inotifyForDir(dir, level, hashflag, rpcflag)
	case daemon:
		// rpcreport("aaa")
		return
	}
	flag.Usage()
}
