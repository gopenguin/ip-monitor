FROM golang:1.10-alpine as builder

RUN apk add --no-cache git build-base

WORKDIR /go/src/github.com/gopenguin/ip-monitor

RUN go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-linkmode external -extldflags \"-static -lc\" -w -s" -o ip-monitor .


FROM scratch
COPY --from=builder /go/src/github.com/gopenguin/ip-monitor/ip-monitor .
CMD ["/ip-monitor"]

