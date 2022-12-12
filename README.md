# WebServer-FinalExam

## prerequisite
* Go 1.19
* Already configured Go development environment

## Instructions

#### First , You must switch to the project root directory , Execute the following command in order to compile the entire project

```@shell
make clean build
``` 

#### This produces executables for many platforms, choose the one that suits you

## usage

#### The 0.0.1 Version is Command Line Tool But Now is a WebServer, There are some endpoints exposed

First you must start the WebServer after you complied the source with thd command 

Linux:
```shell
./WEB-FinalExam_linux_amd64
```
Windows:
```shell
./_windows_amd64.exe
```


## pay attention
* Before using the audio conversion function, please make sure that the operating 
  system environment variable is configured For Example **AWS_ACCESS_KEY_ID** ,**AWS_SECRET_ACCESS_KEY**,
  This project only supports calling Region:ap-northeast-1  AWS services.

* Due to network restrictions in the China region, calling AWS for audio recognition, 
  if used without a proxy, may cause the network to be unable to connect and keep getting stuck.
