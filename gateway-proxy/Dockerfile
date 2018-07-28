FROM golang:latest

WORKDIR /go/src/app
COPY . .
COPY ./pb /go/src/app/pb
COPY ./run.sh /var/run.sh
RUN go get -d -v ./...
RUN go build -v ./...
# ENTRYPOINT ["/go/src/app"]
ADD run.sh /root/start.sh

CMD ["/bin/bash", "/root/start.sh"]

EXPOSE 8080