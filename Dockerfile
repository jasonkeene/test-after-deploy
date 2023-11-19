FROM --platform=$BUILDPLATFORM golang:1.21 as builder

WORKDIR /build

ARG TARGETOS
ARG TARGETARCH

# CMD must be set to the name of the command to build in the cmd/ directory.
ARG CMD

# MODE must be set to "build" or "test". It represents the go build mode. It
# defaults to "build".
ARG MODE=build

ENV GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    CGO_ENABLED=0

COPY go.mod go.sum ./

RUN go mod download

COPY / .

RUN [ "$MODE" = "build" ] || [ "$MODE" = "test" ] || exit 1 && \
    go $MODE \
    -a \
    -installsuffix nocgo \
    $([ $MODE = "test" ] && echo "-c") \
    -o bin \
    ./cmd/$CMD

FROM --platform=$TARGETPLATFORM scratch

WORKDIR /srv
COPY --from=builder /build/bin .
ENTRYPOINT [ "/srv/bin" ]
