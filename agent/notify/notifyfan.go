package notify

import (
	"fmt"
	"log"
	"os"

	// "github.com/google/logger"
	"github.com/s3rj1k/go-fanotify/fanotify"
	"golang.org/x/sys/unix"
)

//RunFanotify 运行Fanotify
func RunFanotify(mountpoint string, hashflag bool, serverip string) {
	// logger.SetFlags(log.Lshortfile)

	notify, err := fanotify.Initialize(
		unix.FAN_CLOEXEC|unix.FAN_CLASS_NOTIF|unix.FAN_UNLIMITED_QUEUE|unix.FAN_UNLIMITED_MARKS,
		os.O_RDONLY|unix.O_LARGEFILE|unix.O_CLOEXEC,
	)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if err = notify.Mark(
		unix.FAN_MARK_ADD|unix.FAN_MARK_MOUNT,
		unix.FAN_MODIFY|unix.FAN_CLOSE_WRITE,
		unix.AT_FDCWD,
		mountpoint,
	); err != nil {
		log.Fatalf("%v\n", err)
	}

	f := func(notify *fanotify.NotifyFD) (string, error) {
		data, err := notify.GetEvent(os.Getpid())
		if err != nil {
			return "", fmt.Errorf("%w", err)
		}

		if data == nil {
			return "", nil
		}

		defer data.Close()

		path, err := data.GetPath()
		if err != nil {
			return "", err
		}

		if data.MatchMask(unix.FAN_CLOSE_WRITE) || data.MatchMask(unix.FAN_MODIFY) {
			return fmt.Sprintf("PID:%d %s", data.GetPID(), path), nil
		}

		return "", fmt.Errorf("fanotify: unknown event")
	}

	for {
		str, err := f(notify)
		if err == nil && len(str) > 0 {
			fmt.Printf("%s\n", str)
		}

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}
