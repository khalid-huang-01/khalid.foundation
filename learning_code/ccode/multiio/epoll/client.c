#include <sys/types.h>
#include <sys/socket.h>
#include <stdio.h>
#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>
#include <sys/select.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

#define MAXLINE 1024
char sendhead[MAXLINE];


int main(int argc , char* argv[])
{
    int sockfd;
    struct sockaddr_in servaddr;
    char *info="cxt";
    int   maxfdp1, stdineof;
    fd_set  rset;
    char  recvbuf[MAXLINE],tmp[128],sendbuf[MAXLINE];
    int   n,len,fd;

    if(argc!=3){
        printf("useage:client address port ");
        exit(0);
    }
    if((sockfd=socket(AF_INET,SOCK_STREAM,0))==-1  )
    {
        perror("socket");
        exit(1);
    }
    printf("%s connect server/n",info);

    bzero(&servaddr,sizeof(servaddr));
    servaddr.sin_family=AF_INET;
    servaddr.sin_port=htons(atoi(argv[2]));
    inet_pton(AF_INET,argv[1],&servaddr.sin_addr);

    if( ( connect(sockfd,(struct sockaddr*)&servaddr,sizeof(servaddr))  )<0)
    {
        perror("connect");
        exit(1);
    }
    send(sockfd,info,strlen(info),0);

    for ( ; ; ) {
        FD_ZERO(&rset);
        FD_SET(sockfd, &rset);
        FD_SET(0, &rset);
        maxfdp1=sockfd+1;

        if((select(maxfdp1, &rset, NULL, NULL, NULL) )<=0){
            perror("select");
        }else{
            if (FD_ISSET(0,&rset)){
                fgets(sendbuf, MAXLINE, stdin);
                n=send(sockfd,sendbuf,strlen(sendbuf)-1,0);
                if(n>0)
                    printf("send: %s",sendbuf);
                else
                    printf("send: %s error,the erro cause is %s:%s/n",sendbuf,errno,strerror(errno));
                bzero(sendbuf,strlen(sendbuf));
            }

            if (FD_ISSET(sockfd, &rset)) { /* socket is readable */
                n=recv(sockfd, recvbuf, MAXLINE,0) ;
                if(n<0) {
                    perror("str_cli: server terminated prematurely");
                }else if(n==0)
                {
                    printf("sever shutdown!");
                    exit(-1);
                }
                //recvbuf
                recvbuf[n]='/0';
                printf("receive :%s/n",recvbuf);
                fflush(stdout);
                bzero(recvbuf,strlen(recvbuf));
            }
        }
    }
    exit(0);
}
