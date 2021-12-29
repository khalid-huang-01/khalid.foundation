//
// Created by Administrator on 2021/12/15.
//

#include <string.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <netdb.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>

int main(int argc, char **argv) {
    int c_fd;
    int pid;
    char sendbuf[128];
    char recvbuf[128];
    struct sockaddr_in c_addr;

    memset(&c_addr, 0, sizeof(struct sockaddr_in));

    // 1. socket
    c_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (c_fd == -1) {
        perror("socket");
        exit(-1);
    }

    // 2. connect
    c_addr.sin_family = AF_INET;
    c_addr.sin_port = htons(atoi(argv[2]));
    inet_aton(argv[1], &c_addr.sin_addr);
    connect(c_fd, (struct sockaddr *)&c_addr, sizeof(struct sockaddr_in));

    printf("=================连接成功================");

    pid = fork();
    // pid > 0 是主进程
    if (pid > 0) {
        while (1) {
            memset(recvbuf, 0, sizeof(recvbuf));
            read(c_fd, recvbuf, 128);
            if(strstr(recvbuf, "client quit") != NULL) {
                printf("==================成功退出==================\n");
                exit(0);
            }
            printf("服务端: %s", recvbuf);
        }
    } else if (pid == 0){
        // 子进程监听输入
        while (1) {
            memset(sendbuf, 0, sizeof(sendbuf));
            fgets(sendbuf, 128, stdin);
            write(c_fd, sendbuf, strlen(sendbuf));
            if(strstr(sendbuf, "quit") != NULL) {
                exit(0);
            }
        }
    }
    close(c_fd);
    return 0;
}
