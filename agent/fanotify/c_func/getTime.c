#include <time.h>
#include <stdio.h>

void getTime()
{
    time_t timep;
    time(&timep);
    printf("%s ", asctime(gmtime(&timep)));
}