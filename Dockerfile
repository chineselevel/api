FROM mischief/docker-golang
ENV HOME /root
RUN mkdir -p $GOPATH/src/github.com/chineselevel
RUN cd $GOPATH/src/github.com/chineselevel/; git clone https://github.com/chineselevel/api.git
RUN cd $GOPATH/src/gitub.com/chineselevel/api; go get -v ./...
WORKDIR /root
CMD go run $GOPATH/src/github.com/chineselevel/api/main.go --port 8000

EXPOSE 8000
