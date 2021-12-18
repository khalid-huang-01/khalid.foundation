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
    if(argc != 3) {
        printf("Usage: %s <addr> <port> \n", argv[0]);
        exit(1);
    }

    int s_fd;
    int c_fd;
    int pid;
    char sendbuf[128];
    char recvbuf[128];

    struct sockaddr_in s_addr;
    struct sockaddr_in c_addr;


    memset(&s_addr, 0, sizeof(struct sockaddr_in));
    memset(&c_addr, 0, sizeof(struct sockaddr_in));

    //1. socket
    s_fd = socket(AF_INET, SOCK_STREAM, 0);
    if(s_fd == -1) {
        perror("socket");
        exit(-1);
    }


    // 2. bind
    s_addr.sin_family = AF_INET;
    s_addr.sin_port = htons(atoi(argv[2]));
    inet_aton(argv[1], &s_addr.sin_addr);
    bind(s_fd, (struct sockaddr *)&s_addr, sizeof(struct sockaddr_in));

    // 3. listen
    listen(s_fd, 10);

    printf("=================listen 成功================");

    // 4. accept
    int c_len = sizeof(struct sockaddr_in);
    while (1) {
        c_fd = accept(s_fd,(struct sockaddr *)&c_addr, &c_len);
        if (c_fd == -1) {
            perror("accept");
            exit(-1);
        }
        printf("===============客户端已连接===========\n");
        pid = fork();
        // pid == 0 表示子进程
        if (pid == 0) {
            while (1) {
                memset(recvbuf, 0, sizeof(recvbuf));
                read(c_fd, recvbuf, 128);
                if (strstr(recvbuf, "quit") != NULL) {
                    write(c_fd, "client quit", 20);
                    printf("===========客户端已退出============\n");
                    break;
                } else {
                    printf("客户端： %s", recvbuf);
                }
            }
        } else if (pid > 0) {
            while (1) {
                memset(sendbuf, 0, sizeof(sendbuf));
                fgets(sendbuf, 128, stdin);
                write(c_fd, sendbuf, strlen(sendbuf));
            }
        }
    }
    close(s_fd);
    close(c_fd);
    return 0;
}