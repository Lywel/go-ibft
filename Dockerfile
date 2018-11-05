FROM golang:latest
WORKDIR /go/src/github.com/Lywel/go-ibft
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
COPY --from=0 /go/src/github.com/Lywel/go-ibft/go-ibft .
ENTRYPOINT ['/go-ibft']

