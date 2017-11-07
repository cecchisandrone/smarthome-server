FROM golang

ADD . /go/src/github.com/cecchisandrone/smarthome-server

RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/cecchisandrone/smarthome-server

RUN dep ensure --vendor-only

RUN go build -o smarthome-server *.go

ENTRYPOINT /go/src/github.com/cecchisandrone/smarthome-server/smarthome-server

EXPOSE 8080