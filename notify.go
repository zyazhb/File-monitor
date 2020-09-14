package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func inotify(filenames []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("[*]event:", event)
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

func inotifyForDir(dir string) {
	filenames, err := GetAllFile(dir)
	if err != nil {
		log.Print("[-]Error: ", err)
	}
	inotify(filenames)
}
