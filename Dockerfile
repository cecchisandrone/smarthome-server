FROM golang:1.15.15

ENV GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0

RUN curl -sL -o /bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /bin/dep

ADD . /go/src/github.com/cecchisandrone/smarthome-server

WORKDIR /go/src/github.com/cecchisandrone/smarthome-server

RUN dep ensure

RUN go build -a -installsuffix cgo -o smarthome-server main.go
