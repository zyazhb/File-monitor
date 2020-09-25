#define _GNU_SOURCE
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

void handle_events(int fd, int mount_fd);

int main(int argc, char *argv[])
{
    int fd, ret, mount_fd, poll_num;
    char buf; // 在判断stdin时使用
    nfds_t nfds;
    struct pollfd fds[2];

    if (argc != 2)
    {
        fprintf(stderr, "Invalid number of command line arguments.\n");
        exit(EXIT_FAILURE);
    }

    mount_fd = open(argv[1], O_DIRECTORY | O_RDONLY);
    if (mount_fd == -1)
    {
        perror(argv[1]);
        exit(EXIT_FAILURE);
    }

    fd = fanotify_init(FAN_CLASS_NOTIF | FAN_REPORT_FID | FAN_NONBLOCK, 0);
    if (fd == -1)
    {
        perror("fanotify_init");
        exit(EXIT_FAILURE);
    }

    ret = fanotify_mark(fd, FAN_MARK_ADD | FAN_MARK_ONLYDIR, FAN_CREATE | FAN_ONDIR | FAN_DELETE, AT_FDCWD, argv[1]);
    if (ret == -1)
    {
        perror("fanotify_mark");
        exit(EXIT_FAILURE);
    }

    // 准备poll
    nfds = 2;

    fds[0].fd = STDIN_FILENO;
    fds[0].events = POLLIN;

    fds[1].fd = fd;
    fds[1].events = POLLIN;

    // 开始监听对应文件
    printf("Listening for events.\n");

    while (1)
    {
        poll_num = poll(fds, nfds, -1); // -1表示无限等待
        if (poll_num == -1)
        {
            if (errno == EINTR) // 遇到错误阻塞的情况
            {
                continue;
            }
            perror("poll");
            exit(EXIT_FAILURE);
        }

        if (poll_num > 0)
        {
            if (fds[0].revents & POLLIN)
            {
                // 处理stdin
                while (read(STDIN_FILENO, &buf, 1) > 0 && buf != '\n')
                {
                    continue;
                }
                break;
            }

            if (fds[1].revents & POLLIN)
            {
                handle_events(fd, mount_fd);
            }
        }
    }

    printf("All events processed successfully. Program exiting.\n");
    exit(EXIT_SUCCESS);
}

void handle_events(int fd, int mount_fd)
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
            perror("read");
            exit(EXIT_FAILURE);
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
                fprintf(stderr, "Received unexpected event info type.\n");
                exit(EXIT_FAILURE);
            }

            if (metadata->mask == FAN_CREATE)
            {
                printf("FAN_CREATE (file created):\n");
            }

            if (metadata->mask == FAN_DELETE)
            {
                printf("FAN_DELETE (file deleted):\n");
            }

            if (metadata->mask == (FAN_CREATE | FAN_ONDIR))
            {
                printf("FAN_CREATE | FAN_ONDIR (subdirectory created):\n");
            }

            if (metadata->mask == (FAN_DELETE | FAN_ONDIR))
            {
                printf("FAN_DELETE | FAN_ONDIR (subdirectory deleted):\n");
            }

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
                    perror("open_by_handle_at");
                    exit(EXIT_FAILURE);
                }
            }

            snprintf(procfd_path, sizeof(procfd_path), "/proc/self/fd/%d", event_fd);

            // 将软连接读成绝对路径
            path_len = readlink(procfd_path, path, sizeof(path) - 1);
            if (path_len == -1)
            {
                perror("readlink");
                exit(EXIT_FAILURE);
            }

            path[path_len] = '\0';
            printf("\tDirectory '%s' has been modified.\n", path);

            // Close associated file descriptor for this event
            close(event_fd);
        }
    }
}