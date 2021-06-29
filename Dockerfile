ARG alpine_version=3.8
ARG go_version=1.16.5

# Server binary builder
FROM golang:${go_version}-alpine${alpine_version} as server_builder

RUN apk add --update --no-cache curl build-base

ENV REPO_DIR ${GOPATH}/src/github.com/angelinahung/product-category
ENV SERVER_DIR ${REPO_DIR}

COPY . ${SERVER_DIR}/

WORKDIR ${SERVER_DIR}

RUN go mod download
RUN go build -o /opt/server

CMD ["/bin/sh"]


# Deployable with server binary and UI dist
FROM alpine:${alpine_version}

WORKDIR /root/
COPY --from=server_builder /opt/server /root/server

EXPOSE 8000

CMD ["./server", "--dev-mode", "--host=0.0.0.0", "--port=8000"]


