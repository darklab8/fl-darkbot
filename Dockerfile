FROM golang:1.23-bullseye as dependencies

RUN apt update
RUN apt install -y build-essential
RUN apt-get install ca-certificates -y
RUN gcc --version

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download -x

FROM dependencies as build

RUN mkdir data
COPY main.go ./
COPY app app
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -v -o main main.go

FROM debian:11.6-slim as runner
WORKDIR /code
RUN mkdir data
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ARG BUILD_VERSION
ENV BUILD_VERSION="${BUILD_VERSION}"
CMD ./main run
