package main

import (
	"crypto/sha256"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func inotify(filenames []string, hashflag bool, rpcflag bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var filehash = new([]byte)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("[*]event:", event)

				if hashflag {
					*filehash = calcHash(event.Name)
				}
				if rpcflag {
					rpcreport(event, *filehash)
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("[*]modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("[-]error:", err)
			}
		}
	}()

	for _, filename := range filenames {
		log.Printf("[+]Add watcher: " + filename)
		err = watcher.Add(filename)
		if err != nil {
			log.Fatal(err)
		}
	}
	<-done

}

func inotifyForDir(dir string, hashflag bool, rpcflag bool) {
	filenames, err := GetAllFile(dir)
	if err != nil {
		log.Print("[-]Error: ", err)
	}
	inotify(filenames, hashflag, rpcflag)
}

func calcHash(filename string) []byte {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Println("文件读取失败！")
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalln(err)
	}
	sum := hash.Sum(nil)

	return sum
}
