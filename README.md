ChineseLevel API
================

The ChineseLevel API is a Go server that provides several Chinese-related functions in one convenient RESTful JSON API.


Quickstart
----------

```shell
$ go run main.go --port 7000
```

Endpoints
----------

### /split [GET/POST]

Takes a Chinese string and returns a tokenized array built out of the words in the string.

 - Parameters:
   + text [string]

##### Example:

Request:
```
/split?text=我叫何曼
```

Response:
```
{
    "text": [
        "我",
        "叫",
        "何曼"
    ]
}
```

*******************************************

### /words [GET/POST]

Takes a Chinese string and returns a tokenized array built out of the words in the string (like /split), together with extra information, like each individual word's rank. The returned rank is -1 if the word was not found.

 - Parameters:
   + text [string]

##### Example:

Request:
```
/split?text=我叫何曼
```

Response:
```
{
    "words": [
        {
            "rank": 7,
            "word": "我"
        },
        {
            "rank": 156,
            "word": "叫"
        },
        {
            "rank": -1,
            "word": "何曼"
        }
    ]
}
```

*******************************************

### /rank [GET/POST]

##### Parameters:

 - Parameters:
   + text [string]

##### Example:

*******************************************

### /analyze [GET/POST]

##### Parameters:

 - Parameters:
   + text [string]

##### Example:

Request:
```
GET /analyze?text=她是一位患有先天小兒麻痺症的媽媽，不論刮風下雨她都每天在碼頭用自己殘疾的手腳來給他的兒子掙取學費
```

Response:
```
{
    "hsk": 6,
    "percentile": {
        "80": 19355,
        "90": 121684,
        "95": 138696,
        "99": 253514
    },
    "score": 100
}
```

Installation
----------

(Work in progress - You don't require Docker to get it running, but I aim to make it as easy as just downloading the Docker box and running the server.)

Install Docker ([instructions](http://docs.docker.io/en/latest/installation/ubuntulinux/)):

```bash
sudo apt-get update
sudo apt-get install linux-image-extra-`uname -r`
sudo sh -c "wget -qO- https://get.docker.io/gpg | apt-key add -"
sudo sh -c "echo deb http://get.docker.io/ubuntu docker main\
> /etc/apt/sources.list.d/docker.list"
sudo apt-get update
sudo apt-get install lxc-docker
```

And run it to confirm it worked (type `exit` to exit):

```bash
sudo docker run -i -t ubuntu /bin/bash
```