cmake_minimum_required(VERSION 3.7)
project(glogwrapper)

set(CMAKE_CXX_STANDARD 14)

add_library(glogwrapper glog_logging.cpp glog_logging.h)
target_link_libraries(glogwrapper glog)
