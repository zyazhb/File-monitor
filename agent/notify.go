package main

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/google/logger"
)

func inotify(filenames []string, hashflag bool, rpcflag bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
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
				logger.Error("[-]error:", err)
			}
		}
	}()

	for _, filename := range filenames {
		logger.Info("\033[1;32m [+]Add watcher: " + filename + " \033[0m")
		err = watcher.Add(filename)
		if err != nil {
			logger.Fatal(err)
		}
	}
	logger.Info("\033[1;30m [*]Add watcher Done! \033[0m")
	<-done

}

// PrintInotifyOp 显示Inotify的Op
func PrintInotifyOp(Name string, Op fsnotify.Op) {
	if Op&fsnotify.Create == fsnotify.Create {
		logger.Info("\033[1;33m [*]Create file:", Name, "\033[0m")
	}
	if Op&fsnotify.Remove == fsnotify.Remove {
		logger.Info("\033[1;33m [*]Remove file:", Name, "\033[0m")
	}
	if Op&fsnotify.Write == fsnotify.Write {
		logger.Info("\033[1;33m [*]Write file:", Name, "\033[0m")
	}
	if Op&fsnotify.Rename == fsnotify.Rename {
		logger.Info("\033[1;33m [*]Rename file:", Name, "\033[0m")
	}
	if Op&fsnotify.Chmod == fsnotify.Chmod {
		logger.Info("\033[1;33m [*]Chmod file:", Name, "\033[0m")
	}
}

func inotifyForDir(dir string, level int, hashflag bool, rpcflag bool) {
	if !strings.HasSuffix(dir, "/"){
		dir += "/"
	}
	filenames, err := GetAllFile(dir, level)
	if err != nil {
		logger.Error("[-]Error: ", err)
	}
	inotify(filenames, hashflag, rpcflag)
}
