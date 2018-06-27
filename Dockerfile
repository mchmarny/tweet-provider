FROM golang:1.10.3
WORKDIR /go/src/github.com/mchmarny/twitter-to-pubsub-event-pump/
COPY . .

# restore to pinnned versions of dependancies
RUN go get github.com/golang/dep/cmd/dep
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -o tpump .


FROM scratch
COPY --from=0 /go/src/github.com/mchmarny/twitter-to-pubsub-event-pump/tpump .
ENTRYPOINT ["/tpump"]