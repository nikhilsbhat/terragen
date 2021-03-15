### Description: Dockerfile for Terragen
FROM alpine:3.11

COPY terragen /

# Starting
ENTRYPOINT [ "/terragen" ]