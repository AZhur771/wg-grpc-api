FROM golang:1.18.8 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o /opt/wg-grpc-api .

FROM alpine:3.18.3

COPY migrations "/migrations/"
COPY --from=build "/opt/wg-grpc-api" "/opt/wg-grpc-api"

CMD "/opt/wg-grpc-api"
