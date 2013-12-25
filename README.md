# ChineseLevel API

The ChineseLevel API is a Go server that provides several Chinese-related functions in one convenient RESTful JSON API.


## Quickstart


```shell
$ go run main.go --port 7000
```

## Endpoints

#### /split [GET/POST]

Takes a Chinese string and returns a tokenized array built out of the words in the string.

##### Parameters: 

 - text [string]

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

#### /rank [GET/POST]

Parameters: 

Example: 

## Installation

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