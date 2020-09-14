package main

import (
	"flag"
	"log"
)

//参数名
var (
	h bool

	f   string
	dir string
	daemon bool
)

//初始化参数
func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&f, "f", "", "choose a file to monitor")
	flag.StringVar(&dir, "dir", "", "choose a dir to monitor")
	flag.BoolVar(&daemon, "daemon", false, "Start in daemon mode")
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
		inotify(filename)

	case len(dir) != 0:
		log.Print("[+]Start dirwalk: " + dir)
		inotifyForDir(dir)
	case daemon:
		rpcreport()
	}
	flag.Usage()
}
