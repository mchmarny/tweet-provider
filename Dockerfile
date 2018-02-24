FROM golang:1.9.2
WORKDIR /go/src/github.com/mchmarny/twitter-to-pubsub-event-pump/
COPY . .
RUN go get github.com/tools/godep
RUN godep restore
RUN CGO_ENABLED=0 GOOS=linux go build -o tpump .

FROM scratch
COPY --from=0 /go/src/github.com/mchmarny/tpump .
ENTRYPOINT ["/tpump"]
