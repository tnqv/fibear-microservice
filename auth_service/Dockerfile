FROM golang:latest

WORKDIR /go/src/app
COPY . .
COPY ./pb /go/src/app/pb
COPY ./common /go/src/app/common
RUN go get -d -v ./...
RUN go build -v ./...
# ENTRYPOINT /go/src/app
# ENTRYPOINT ["/go/src/app"]
# CMD ["/var/run.sh"]
ADD run.sh /root/start.sh

CMD ["/bin/bash", "/root/start.sh"]

EXPOSE 7800