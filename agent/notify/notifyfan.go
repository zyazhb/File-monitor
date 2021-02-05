package notify

import (
	"fmt"
	"os"

	"github.com/google/logger"
	"github.com/s3rj1k/go-fanotify/fanotify"
	"golang.org/x/sys/unix"
)

//RunFanotify 运行Fanotify监控 
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

	for {
		event, pid, path, err := RealFanotify(notify)
		var filehash string
		if hashflag {
			filehash = calcHash(path) //计算hash
		}
		if serverip != "" {
			go rpcreportfan(path, filehash, event, serverip) //rpc上报
		}
		if err == nil && pid != -1 {
			logger.Info("\033[1;33m [*]PID:", pid, " ", event, " ", path, "\033[0m")
		}

		if err != nil {
			logger.Error("\033[1;31m [-]error:", err, "\033[0m")
		}
	}
}

//RealFanotify 真正的fanotify处理函数 (动作，PID，路径，错误)
func RealFanotify(notify *fanotify.NotifyFD) (string, int, string, error) {
	data, err := notify.GetEvent(os.Getpid())
	if err != nil {
		return "", -1, "", fmt.Errorf("%w", err)
	}

	if data == nil {
		return "", -1, "", err
	}

	defer data.Close()

	path, err := data.GetPath()
	if err != nil {
		return "", -1, "", err
	}

	switch {
	case data.MatchMask(unix.FAN_ACCESS):
		return "FAN_ACCESS", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_OPEN):
		return "Open", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_ATTRIB):
		return "FAN_ATTRIB", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_CREATE):
		return "FAN_CREATE", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_DELETE):
		return "FAN_DELETE", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_MODIFY):
		return "FAN_MODIFY", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_CLOSE):
		return "FAN_CLOSE", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_MOVE):
		return "FAN_MOVE", data.GetPID(), path, nil
	case data.MatchMask(unix.FAN_CLOSE_WRITE):
		return "FAN_CLOSE_WRITE", data.GetPID(), path, nil
	}

	return "", -1, "", fmt.Errorf("fanotify: unknown event")

}
