#include <stdio.h>
#include <stdlib.h>
#include <sys/resource.h>

void ulimit(void)
{
    char *val = getenv("_ULIMIT");
    if (val == NULL || *val != '1')
    {
        return;
    }

    struct rlimit rlim;
    int ret;

    ret = getrlimit(RLIMIT_NOFILE, &rlim);
    if (ret == -1)
    {
        perror("getrlimit");
        return;
    }

    printf(" child: rlim_cur: %ld rlim_max: %ld \n", rlim.rlim_cur, rlim.rlim_max);
    exit(0);
}