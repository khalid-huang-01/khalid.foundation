cmake_minimum_required(VERSION 3.10)
project(multiio C)

set(CMAKE_C_STANDARD 99)

#add_executable(multiio select/server.c)
add_executable(multiio epoll/server.c)
