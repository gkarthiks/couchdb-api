FROM golang:1.9 AS builder
RUN go version

COPY . "/go/src/github.com/gkarthiks/couchdb-api"
WORKDIR "/go/src/github.com/gkarthiks/couchdb-api"

#RUN go get -v -t  .
RUN set -x && \
    go get github.com/sirupsen/logrus && \  
    go get github.com/leesper/couchdb-golang

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /couchdb-api


FROM scratch

COPY --from=builder /couchdb-api .
ENTRYPOINT [ "/couchdb-api"]