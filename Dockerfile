# Build image for apps that depend on latlon2h3 / uber/h3-go.
# The libh3-dev package provides the C headers and shared library
# required by the CGo bindings at compile time.
#
# For a runtime-only image, install libh3-4 (not -dev) and
# copy the compiled binary from a build stage.

FROM golang:1.22-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends \
    libh3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
