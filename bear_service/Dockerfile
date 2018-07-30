FROM golang:latest

WORKDIR /go/src/app
COPY . .
COPY ./pb /go/src/app/pb
COPY ./common /go/src/app/common
RUN go get -d -v ./...
# RUN rm -rf /go/src/github.com/appleboy/gorush/rpc/proto
RUN rm -rf /go/src/github.com/appleboy/gorush/vendor/google.golang.org/grpc
RUN go build -v ./...

ADD run.sh /root/start.sh

CMD ["/bin/bash", "/root/start.sh"]

EXPOSE 7801