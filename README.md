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

#### First , This is a command line tool ，Here are his parameters

>  -t < int >  Specify target AlbumId  Default **-1**
> 
>  -f < string > Specify the target AlbumId file Default **""**
> 
>  -d < bool >  Whether to disable downloading audio on execution Default **false**
> 
>  -c < int > Specify The amount of audio that will be converted Default **0**
> 
>  -s < string > Specify the storage parent Dir Path Default **" ./ "**
> 
>  -n < bool > Whether to convert audio only Default **false**

#### use Example: In Linux

```shell
./WEB-FinalExam_linux_amd64 -t=20403413 
```

## pay attention
* Before using the audio conversion function, please make sure that the operating 
  system environment variable is configured For Example **AWS_ACCESS_KEY_ID** ,**AWS_SECRET_ACCESS_KEY**,
  This project only supports calling Region:ap-northeast-1  AWS services.

* Due to network restrictions in the China region, calling AWS for audio recognition, 
  if used without a proxy, may cause the network to be unable to connect and keep getting stuck.
