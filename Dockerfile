FROM golang:1.21-bullseye

RUN useradd --create-home --home-dir /home/test test
ENV USER test

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download

COPY main.go lib.go cidr.go ./
RUN go build -o /usr/bin/cidr

USER test
WORKDIR /home/test

ENTRYPOINT ["cidr"]

