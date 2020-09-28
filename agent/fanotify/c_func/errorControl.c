#include <errno.h>
#include <stdio.h>
#include "errorMacro.h"

void handleError(unsigned int err, char *file)
{
    switch (err)
    {
    case ERROR_WRONG_INPUT:
        fprintf(stderr, "Invalid number of command line arguments.\n");
        break;
    case ERROR_OPEN:
        perror(file);
        break;
    case ERROR_FANOTIFY_INIT:
        perror("fanotify_init");
        break;
    case ERROR_FANOTIFY_MARK:
        perror("fanotify_mark");
        break;
    case ERROR_POLL:
        perror("poll");
        break;
    case ERROR_READ:
        perror("read");
        break;
    case ERROR_EVENT_INFO_TYPE:
        fprintf(stderr, "Received unexpected event info type.\n");
        break;
    case ERROR_OLD_EVENT_FD:
        perror("open_by_handle_at");
        break;
    case ERROR_READLINK:
        perror("readlink");
        break;
    default:
        printf("no error!!!");
        break;
    }
}
