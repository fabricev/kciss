FROM ghcr.io/aquasecurity/trivy:0.18.3 as trivy

FROM golang:1.16 as builder
WORKDIR /go/src/github.com/fabricev/kciss/cmd
COPY . /go/src/github.com/fabricev/kciss/
RUN CGO_ENABLED=0 GOOS=linux go build -o ../kciss

FROM alpine:3.14.0
LABEL org.opencontainers.image.source https://github.com/fabricev/kciss
RUN apk --no-cache add ca-certificates=20191127-r5
WORKDIR /root/
COPY --from=trivy /usr/local/bin/trivy /usr/local/bin/trivy
COPY --from=builder /go/src/github.com/fabricev/kciss/kciss /usr/local/bin/kciss
CMD ["/usr/local/bin/kciss"]
