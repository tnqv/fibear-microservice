FROM golang:1.7

WORKDIR /go/src/app
COPY . .
COPY ./pb /go/src/app/pb
COPY ./common /go/src/app/common
RUN go get -d -v ./...
RUN go build -v ./...

EXPOSE 7800