FROM golang:1.17.4-alpine AS builder
ENV GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0
ADD . /go/src/github.com/cecchisandrone/smarthome-server
WORKDIR /go/src/github.com/cecchisandrone/smarthome-server
RUN go mod download
RUN go build -a -installsuffix cgo -o smarthome-server main.go
FROM resin/raspberry-pi-alpine:3.7
COPY --from=builder /go/src/github.com/cecchisandrone/smarthome-server/smarthome-server ./
RUN apk add --no-cache tzdata
ADD config /config
CMD ["./smarthome-server"]
EXPOSE 8080