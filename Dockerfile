FROM golang

ENV GOOS=linux GOARCH=arm GOARM=6

RUN curl -sL -o /bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /bin/dep

ADD . /go/src/github.com/cecchisandrone/smarthome-server

WORKDIR /go/src/github.com/cecchisandrone/smarthome-server

RUN dep ensure

RUN go build -o /bin/smarthome-server *.go

RUN rm -rf vendor && rm -rf /go/src

ENTRYPOINT /bin/smarthome-server

EXPOSE 8080