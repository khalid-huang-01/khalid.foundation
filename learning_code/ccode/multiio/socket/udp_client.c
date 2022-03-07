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

int main(int argc, char* argv[])
{
    int c_fd;
    char sendbuf[128];
    char recvbuf[128];
    socklen_t len;
    struct sockaddr_in s_addr,c_addr;

    //1.socket
    c_fd = socket(AF_INET, SOCK_DGRAM, 0);

    memset(&s_addr, 0, sizeof(s_addr));
    s_addr.sin_family = AF_INET;
    s_addr.sin_port = htons(atoi(argv[2]));
    inet_aton(argv[1],&s_addr.sin_addr);
    printf("==============连接成功==============\n");
    len = sizeof(struct sockaddr_in);

    while(1){
        //2.sendto
        memset(sendbuf, 0, sizeof(sendbuf));
        fgets(sendbuf,128,stdin);
        sendto(c_fd, sendbuf, strlen(sendbuf), 0, (struct sockaddr*)&s_addr, len);
        if(strstr(sendbuf,"quit") != NULL){
            printf("==============成功退出==============\n");
            break;
        }

        //3.recvfrom
        memset(recvbuf, 0, sizeof(recvbuf));
        recvfrom(c_fd, recvbuf, 128, 0, (struct sockaddr*)&c_addr, &len);
        printf("server:%s",recvbuf);
    }

    close(c_fd);
    return 0;
}