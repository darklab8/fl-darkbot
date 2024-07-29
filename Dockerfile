FROM golang:1.22-bullseye as build

RUN apt update
RUN apt install -y build-essential
RUN apt-get install ca-certificates -y
RUN gcc --version

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download -x

RUN mkdir data
COPY main.go ./
COPY app app
RUN go build -v -o main main.go

FROM debian:11.6-slim as runner
WORKDIR /code
RUN mkdir data
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ./main run
