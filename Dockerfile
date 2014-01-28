FROM mischief/docker-golang
ENV HOME /root
RUN mkdir -p $GOPATH/src/github.com/chineselevel
RUN cd $GOPATH/src/github.com/chineselevel/; git clone https://github.com/chineselevel/api.git
RUN cd $GOPATH/src/github.com/chineselevel/api; go get -v ./...
RUN mkdir $GOPATH/src/github.com/chineselevel/api/data
RUN cd $GOPATH/src/github.com/chineselevel/api/data; wget https://dl.dropboxusercontent.com/u/16728281/chineselevel/dict.compressed.tab

CMD cd $GOPATH/src/github.com/chineselevel/api; go run main.go --port 8000
