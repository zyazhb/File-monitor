package main

// #include "./c_func/fanotify_handle.c"
import "C" 
import (
    "fmt"
    "os"
    "unsafe"
    "syscall"
)



func main() {
    var nfds C.nfds_t
    var fd, ret, mount_fd C.int // poll_num careate by poll()
    var buf C.char
    var fds [2]C.struct_pollfd

    argc := len(os.Args)
    argv := C.CString(os.Args[1])

    if argc != 2 {
      fmt.Fprintf(os.Stderr, "Invalid number of command line arguments.\n")
      os.Exit(C.EXIT_FAILURE)
    }

    mount_fd = C.__open(argv, C.O_DIRECTORY | C.O_RDONLY)
    if mount_fd == -1 {
      C.perror(argv)
      os.Exit(C.EXIT_FAILURE)
    }

    fd = C.fanotify_init(C.FAN_CLASS_NOTIF | C.FAN_REPORT_FID | C.FAN_NONBLOCK, 0)
    if fd == -1 {
      C.perror(C.CString("fanotify_init"))
      os.Exit(C.EXIT_FAILURE)
    }

    ret = C.fanotify_mark(fd, C.FAN_MARK_ADD | C.FAN_MARK_ONLYDIR, C.FAN_CREATE | C.FAN_ONDIR | C.FAN_DELETE | C.FAN_MOVE_SELF | C.FAN_MOVE, C.AT_FDCWD, argv)
    if ret == -1 {
      C.perror(C.CString("fanotify_mark"))
      os.Exit(C.EXIT_FAILURE)
    }

    nfds = 2

    fds[0].fd = C.STDIN_FILENO
    fds[0].events = C.POLLIN

    fds[1].fd = fd
    fds[1].events = C.POLLIN

    fmt.Print("Listening for events.\n")

    for {
      poll_num, err_poll := C.poll((*C.struct_pollfd)(unsafe.Pointer(&fds)), nfds, -1)
      if poll_num == -1 {
        if err_poll == syscall.EINTR {
          continue
        }
        C.perror(C.CString("poll"))
        os.Exit(C.EXIT_FAILURE)
      }

      if poll_num > 0 {
        if (fds[0].revents & C.POLLIN) != 0 {
          for ;C.read(C.STDIN_FILENO, unsafe.Pointer(&buf), 1) > 0 && buf != '\n'; {
            continue
          }
          break
        }

        if (fds[1].revents & C.POLLIN) != 0 {
          C.getTime()
          C.handle_events(fd, mount_fd)
        }
      }
    }
    fmt.Print("All events processed successfully. Program exiting.\n")
    os.Exit(C.EXIT_SUCCESS)
}