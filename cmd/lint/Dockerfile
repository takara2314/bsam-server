FROM golang:1.23-alpine

RUN apk add --no-cache git curl && \
    rm -rf /var/cache/apk/*

ENV PATH=$PATH:/go/bin

ENV GOLANGCI_LINT_VERSION=v1.60.1
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

WORKDIR /app
COPY . .

ENTRYPOINT ["golangci-lint"]

CMD ["--version"]
