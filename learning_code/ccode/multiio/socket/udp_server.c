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
    int s_fd;
    char sendbuf[128];
    char recvbuf[128];
    socklen_t len;
    struct sockaddr_in s_addr,c_addr;

    //1.socket
    s_fd = socket(AF_INET, SOCK_DGRAM, 0);
    if(s_fd == -1){
        perror("socket");
        exit(-1);
    }

    //2.bind
    memset(&s_addr, 0, sizeof(s_addr));
    s_addr.sin_family = AF_INET;
    s_addr.sin_port = htons(atoi(argv[2]));
    inet_aton(argv[1],&s_addr.sin_addr);
    bind(s_fd, (struct sockaddr*)&s_addr, sizeof(s_addr));

    len = sizeof(struct sockaddr_in);
    while(1){

        //3.recvfrom
        memset(recvbuf, 0, sizeof(recvbuf));
        recvfrom(s_fd, recvbuf, 128, 0, (struct sockaddr*)&c_addr, &len);
        if(strstr(recvbuf,"quit") != NULL){
            printf("==============客户端已退出==============\n");
        }else{
            printf("client:%s",recvbuf);
        }

        //4.sendto
        memset(sendbuf, 0, sizeof(sendbuf));
        fgets(sendbuf,128,stdin);
        if(strstr(sendbuf,"quit") != NULL){
            printf("==============服务端已退出==============\n");
            break;
        }
        sendto(s_fd, sendbuf, strlen(sendbuf), 0, (struct sockaddr*)&c_addr, len);
    }

    close(s_fd);
    return 0;
}