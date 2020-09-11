package main

import (
	// "fmt"
	"flag"
	"log"
)

//参数名
var (
	h bool

    filename string
	dir string
)

//初始化参数
func init() {
    flag.BoolVar(&h, "h", false, "this help")

    // 注意 `signal`。默认是 -s string，有了 `signal` 之后，变为 -s signal
    flag.StringVar(&filename, "filename", "", "choose a `file` to monitor")
}

func main(){
	flag.Parse()

	log.Print("Start[+]")


	if  len(filename) != 0 {
		inotify(filename)
	}
    if h {
        flag.Usage()
    }

}
	