package main

import (
	"log"

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

				PrintInotifyOp(event.Name, event.Op)

				if hashflag {
					*filehash = calcHash(event.Name)
				}
				if rpcflag {
					go rpcreport(event, *filehash)
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

// PrintInotifyOp 显示Inotify的Op
func PrintInotifyOp(Name string, Op fsnotify.Op){
	if Op&fsnotify.Create == fsnotify.Create {
					log.Println("[*]Create file:", Name)
				}
				if Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("[*]Remove file:", Name)
				}
				if Op&fsnotify.Write == fsnotify.Write {
					log.Println("[*]Write file:", Name)
				}
				if Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("[*]Rename file:", Name)
				}
				if Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("[*]Chmod file:", Name)
				}
}

func inotifyForDir(dir string, hashflag bool, rpcflag bool) {
	filenames, err := GetAllFile(dir)
	if err != nil {
		log.Print("[-]Error: ", err)
	}
	inotify(filenames, hashflag, rpcflag)
}
