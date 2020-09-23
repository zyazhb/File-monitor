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
	daemon   bool
	rpcflag  bool
	hashflag bool
)

//初始化参数
func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&f, "f", "", "choose a file to watch")
	flag.StringVar(&dir, "dir", "", "choose a dir to monitor")
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
		inotifyForDir(dir, hashflag, rpcflag)
	case daemon:
		// rpcreport("aaa")
		return
	}
	flag.Usage()
}
