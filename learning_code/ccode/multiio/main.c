#include <stdio.h>
#include <string.h>
#include <sys/select.h>
#include <arpa/inet.h>
#include <unistd.h>

#define STDIN 0
#define BUF_SIZE 30

// https://blog.csdn.net/y396397735/article/details/55004775
_Noreturn void selectFile() {
    fd_set reads, temps;
    int result, str_len;
    char buf[BUF_SIZE];
    struct timeval timeout;
    FD_ZERO(&reads);
    // 监视文件描述符0的变化，也就是标准输入的变化
    FD_SET(STDIN, &reads);
    /*
    超时不能在此设置！
    因为调用select后，结构体timeval的成员tv_sec和tv_usec的值将被替换为超时前剩余时间.
    调用select函数前，每次都需要初始化timeval结构体变量.
    timeout.tv_sec = 5;
    timeout.tv_usec = 5000;
    */
    while(1) {
        /*将准备好的fd_set变量reads的内容复制到temps变量，因为调用select函数后，除了发生变化的fd对应位外，剩下的所有位
        都将初始化为0，为了记住初始值，必须经过这种复制过程。
        */
        temps = reads;
        // 设置超时
        timeout.tv_sec = 5;
        timeout.tv_usec = 0;

        //调用select函数. 若有控制台输入数据，则返回大于0的整数，如果没有输入数据而引发超时，返回0.
        result = select(1, &temps, 0, 0, &timeout);
        switch (result) {
            case -1:
                perror("select function error");
                break;
            case 0:
                puts("timeout");
            default:
                // 判断STDIN是否有读事件发生
                if(FD_ISSET(STDIN, &temps)) {
                    str_len = read(STDIN, buf, BUF_SIZE);
                    buf[str_len] = 0;
                    printf("message from console: %s", buf);
                }
        }
    }

}

int selectTCP(int argc, char *argv[]) {
    int n, flag;
    char buf[1024];
    fd_set fds;
    struct timeval tv;
    struct sockaddr_in  client_addr;
    int sockfd, fd, sin_size;
    int port;
    struct sockaddr_in addr;
    int fdstd = STDIN;



}

int main(int argc, char *argv[]) {
    selectFile();
}
