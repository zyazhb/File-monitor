package notify

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/google/logger"
)

//RunInotify 运行Inotify监控
func RunInotify(filenames []string, hashflag bool, serverip string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("\033[1;31m", err, "\033[0m")
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

				PrintInotifyOp(event.Name, event.Op) //agent回显

				var filehash string
				if hashflag {
					filehash = calcHash(event.Name) //计算hash
				}
				if serverip != "" {
					go rpcreporti(event, filehash, serverip) //rpc上报
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("\033[1;31m [-]error:", err, "\033[0m")
			}
		}
	}()

	for _, filename := range filenames {
		logger.Info("\033[1;32m [+]Add watcher: " + filename + " \033[0m")
		err = watcher.Add(filename)
		if err != nil {
			logger.Fatal("\033[1;31m", err, "\033[0m")
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

//RunInotifyForDir 递归添加Inotify
func RunInotifyForDir(dir string, level int, hashflag bool, serverip string) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	filenames, err := GetAllFile(dir, level)
	if err != nil {
		logger.Error("\033[1;31m [-]Error: ", err, " \033[0m")
	}
	RunInotify(filenames, hashflag, serverip)
}
