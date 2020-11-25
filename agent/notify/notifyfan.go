package notify

import (
	"fmt"
	"os"

	"github.com/google/logger"
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
		logger.Fatal("\033[1;31m", err, "\033[0m")
	}
	//https://man7.org/linux/man-pages/man2/fanotify_mark.2.html 部分特性在kernel 5.0+才有效
	if err = notify.Mark(
		unix.FAN_MARK_ADD|unix.FAN_MARK_MOUNT,
		unix.FAN_ACCESS|unix.FAN_MODIFY|unix.FAN_OPEN|unix.FAN_CLOSE,
		unix.AT_FDCWD,
		mountpoint,
	); err != nil {
		logger.Fatal("\033[1;31m", err, "\033[0m")
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

		switch {
		case data.MatchMask(unix.FAN_ACCESS):
			return fmt.Sprintf("FAN_ACCESS file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_OPEN):
			return fmt.Sprintf("Open file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_ATTRIB):
			return fmt.Sprintf("FAN_ATTRIB file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_CREATE):
			return fmt.Sprintf("FAN_CREATE file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_DELETE):
			return fmt.Sprintf("FAN_DELETE file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_MODIFY):
			return fmt.Sprintf("FAN_MODIFY file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_CLOSE):
			return fmt.Sprintf("FAN_CLOSE file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_MOVE):
			return fmt.Sprintf("FAN_MOVE file: PID:%d %s", data.GetPID(), path), nil
		case data.MatchMask(unix.FAN_CLOSE_WRITE):
			return fmt.Sprintf("FAN_CLOSE_WRITE file: PID:%d %s", data.GetPID(), path), nil
		}

		return "", fmt.Errorf("fanotify: unknown event")
	}

	for {
		str, err := f(notify)
		if err == nil && len(str) > 0 {
			logger.Info("\033[1;33m [*]", str, "\033[0m")
		}

		if err != nil {
			logger.Error("\033[1;31m [-]error:", err, "\033[0m")
		}
	}
}
