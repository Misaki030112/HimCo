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
./WEB-FinalExam_windows_amd64.exe
```
Second Some of the endpoints we provide

|   endpoint   | request method |         request parameters         |                                                                                                                                  Functional description                                                                                                                                   | 
|:------------:|:--------------:|:----------------------------------:|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|
|    /album    |      GET       | ?id=< INT >&audioDownload=< BOOL > |                                                                        Crawl the Album on the Himalayan website according to the AlbumId, and the parameter audioDownload determines whether to download the audio                                                                        |
|   /convert   |      GET       |     ?id=< INT >&count=< INT >      | Convert the audio downloaded from the endpoint /album to audio-to-text, pass in the AlbumID, and distinguish the album. Note that the /album endpoint must be called before the audio has been downloaded. The parameter count determines the number of converted audio. The default is 1 |
| /analyzeJson |      GET       |           no parameters            |                                                      Analyze the book data on the entire Himalayan website, and form an analysis json report. The report information includes, book vip rate, completion rate, top3 of each category                                                      |

## pay attention
* Before using the audio conversion function, please make sure that the operating 
  system environment variable is configured For Example **AWS_ACCESS_KEY_ID** ,**AWS_SECRET_ACCESS_KEY**,
  This project only supports calling Region:ap-northeast-1  AWS services.

* Due to network restrictions in the China region, calling AWS for audio recognition, 
  if used without a proxy, may cause the network to be unable to connect and keep getting stuck.

* When calling the /analyzeJson endpoint, due to the large amount of data in the entire Himalayan 
  website, the analysis time will be longer, which mainly depends on the speed of your network IO.
