FROM golang:1.19.3-bullseye as build

RUN apt update
RUN apt install -y build-essential
RUN gcc --version

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download -x

RUN mkdir data
COPY main.go ./
COPY cmd cmd
COPY configurator configurator
COPY consoler consoler
COPY discorder discorder
COPY listener listener
COPY scrappy scrappy
COPY settings settings
COPY utils utils
COPY viewer viewer
RUN go build -v -o main main.go

FROM debian:11.6-slim as runner
WORKDIR /code
COPY --from=build /code/main main
CMD ./main run