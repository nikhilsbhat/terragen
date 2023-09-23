### Description: Dockerfile for Terragen
FROM golang:1.21.1-alpine3.18 as builder

RUN go install golang.org/x/tools/cmd/goimports@latest && \
    go install mvdan.cc/gofumpt@latest && \
    go install github.com/daixiang0/gci@latest && \
    wget -q https://github.com/hashicorp/terraform-plugin-docs/releases/download/v0.16.0/tfplugindocs_0.16.0_linux_amd64.zip -O /tmp/tfplugindocs.zip && \
    unzip /tmp/tfplugindocs.zip -d /go/bin

FROM golang:1.21.1-alpine3.18

COPY --from=builder /go/bin/goimports /go/bin/goimports
COPY --from=builder /go/bin/gofumpt /go/bin/gofumpt
COPY --from=builder /go/bin/gci /go/bin/gci
COPY --from=builder /go/bin/tfplugindocs /go/bin/tfplugindocs

COPY terragen /

# Starting
ENTRYPOINT [ "/terragen" ]