FROM golang:1.25-bookworm as dependencies

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
ARG BUILD_VERSION
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags "-X github.com/prometheus/common/version.Version=${BUILD_VERSION}" -v -o main main.go

FROM debian:12.11-slim as runner
WORKDIR /code
RUN mkdir data
ARG BUILD_VERSION
ENV BUILD_VERSION="${BUILD_VERSION}"
COPY --from=build /code/main main
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8000
CMD ./main run
