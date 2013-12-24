ChineseLevel
============

Installation
------------

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