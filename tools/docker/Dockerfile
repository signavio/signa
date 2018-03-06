FROM golang:1.10 as builder

WORKDIR /go/src/github.com/signavio/signa
ADD . .
RUN CGO_ENABLED=0 go build -o /signa -a -tags netgo -ldflags '-w' cmd/signa/main.go

FROM alpine
COPY --from=builder /signa /signa

RUN apk update && apk add curl
RUN curl -LO \
    https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x kubectl && mv kubectl /usr/bin/

ENTRYPOINT ["/signa"]
CMD ["-config", "$SIGNA_CONFIG"]
