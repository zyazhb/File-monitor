package module

// #define _GNU_SOURCE
// #include <errno.h>
// #include <fcntl.h>
// #include <limits.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <sys/types.h>
// #include <sys/stat.h>
// #include <sys/fanotify.h>
// #include <unistd.h>
// #include <poll.h>
// #include <time.h>
import "C" 
import (
    "fmt"
	"os"
	"io"
	"bufio"
    "unsafe"
)

type NotifyFD struct {
	Fd   int
	File *os.File
	Rd io.Reader
}

func initFanotify(fanotifyFlags uint, openFlags uint) (*NotifyFD, error) {
	fd, err := C.fanotify_init(C.uint(fanotifyFlags), C.uint(openFlags)) // 如果是C函数中有errno则第二个返回值即是error类
	if fd == -1 {
		return nil, fmt.Errorf("fanotify: init error %w", err)
	}
	file := os.NewFile(uintptr(fd), "")
	rd := bufio.NewReader(file)

	return &NotifyFD{
		Fd: int(fd),
		File: file,
		Rd: rd,
	}, nil
}

func Mark(handle *NotifyFD, flags uint, mask uint64, dirFd int, path string) error {
	ret, err := C.fanotify_mark(C.int(handle.Fd), C.uint(flags), C.uint64_t(mask), C.int(dirFd), C.CString(path)) 
	if ret == -1 {
		return fmt.Errorf("fanotify: mark error %w", err)
	}
	return nil
}



// func initPoll(count uint32, pipeId ...chan string) (int, error) {
func initPoll(fd int) (int, error) {
	var fds [2]C.struct_pollfd

	fds[0].fd = C.STDIN_FILENO
    fds[0].events = C.POLLIN

    fds[1].fd = C.int(fd)
    fds[1].events = C.POLLIN

	pollNum, err := C.poll((*C.struct_pollfd)(unsafe.Pointer(&fds)), 2, -1)
	if pollNum == -1 {
		return -1, fmt.Errorf("fanotify: mark error %w", err)
	}
	return pollNum, nil
}