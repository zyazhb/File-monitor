#define _GNU_SOURCE
#include "errorMacro.h"
#include <errno.h>
#include <fcntl.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/fanotify.h>
#include <unistd.h>
#include <poll.h>

#define BUF_SIZE 256

int __open(const char *__file, int __oflag)
{
    return open(__file, __oflag);
}

unsigned int handle_events(int fd, int mount_fd)
{
    int event_fd;
    ssize_t len, path_len;
    char path[PATH_MAX];
    char procfd_path[PATH_MAX];
    char events_buf[BUF_SIZE];
    struct file_handle *file_handle;
    struct fanotify_event_metadata *metadata;
    struct fanotify_event_info_fid *fid;

    while (1)
    {
        len = read(fd, (void *)&events_buf, sizeof(events_buf));
        if (len == -1 && errno != EAGAIN) // EAGAIN 表示try again
        {
            return ERROR_READ;
            // perror("read");
            // exit(EXIT_FAILURE);
        }

        if (len < 0)
        {
            break;
        }

        for (metadata = (struct fanotify_event_metadata *)events_buf; FAN_EVENT_OK(metadata, len); FAN_EVENT_NEXT(metadata, len))
        {
            fid = (struct fanotify_event_info_fid *)(metadata + 1);
            file_handle = (struct file_handle *)fid->handle;

            // 需要验证event类型
            if (fid->hdr.info_type != FAN_EVENT_INFO_TYPE_FID)
            {
                return ERROR_EVENT_INFO_TYPE;
                // fprintf(stderr, "Received unexpected event info type.\n");
                // exit(EXIT_FAILURE);
            }

            if (metadata->mask == FAN_CREATE)
            {
                printf("FAN_CREATE (file created):\n");
            }

            if (metadata->mask == FAN_MOVE)
            {
                printf("FAN_MOVE (file renamed):\n");
            }

            if (metadata->mask == FAN_DELETE)
            {
                printf("FAN_DELETE (file deleted):\n");
            }

            if (metadata->mask == (FAN_MOVE | FAN_ONDIR))
            {
                printf("FAN_MOVE | FAN_ONDIR (subdirectory renamed):\n");
            }

            if (metadata->mask == (FAN_CREATE | FAN_ONDIR))
            {
                printf("FAN_CREATE | FAN_ONDIR (subdirectory created):\n");
            }

            if (metadata->mask == (FAN_DELETE | FAN_ONDIR))
            {
                printf("FAN_DELETE | FAN_ONDIR (subdirectory deleted):\n");
            }

            fflush(stdout); // 冲一下缓冲区

            event_fd = open_by_handle_at(mount_fd, file_handle, O_RDONLY); // 判断文件句柄是否被系统删除
            if (event_fd == -1)
            {
                if (errno == ESTALE) // 旧的文件句柄
                {
                    printf("File handle is no longer valid. "
                           "File has been deleted\n");
                    continue;
                }
                else
                {
                    return ERROR_OLD_EVENT_FD;
                    // perror("open_by_handle_at");
                    // exit(EXIT_FAILURE);
                }
            }

            snprintf(procfd_path, sizeof(procfd_path), "/proc/self/fd/%d", event_fd);

            // 将软连接读成绝对路径
            path_len = readlink(procfd_path, path, sizeof(path) - 1);
            if (path_len == -1)
            {
                return ERROR_READLINK;
                // perror("readlink");
                // exit(EXIT_FAILURE);
            }

            path[path_len] = '\0';
            printf("\tDirectory '%s' has been modified.\n", path);
            fflush(stdout);

            // Close associated file descriptor for this event
            close(event_fd);
        }
    }
}