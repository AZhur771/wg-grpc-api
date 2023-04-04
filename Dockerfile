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
        -o /opt/wg-grpc-api ${CODE_DIR}/cmd/wg-grpc-api

FROM alpine:3.9

ENV BIN_FILE "/opt/wg-grpc-api"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
