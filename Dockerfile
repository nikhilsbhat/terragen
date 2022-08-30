### Description: Dockerfile for Terragen
FROM golang:alpine3.15 as builder

RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install mvdan.cc/gofumpt@latest

FROM golang:alpine3.15

COPY --from=builder /go/bin/goimports /go/bin/goimports
COPY --from=builder /go/bin/gofumpt /go/bin/gofumpt
COPY terragen /

# Starting
ENTRYPOINT [ "/terragen" ]