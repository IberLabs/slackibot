FROM golang:1.11.1-alpine AS builder

LABEL maintainer="Manuel Boira <manuelbcd@gmail.com>"

RUN apk update && apk add --no-cache git gcc musl-dev
WORKDIR $GOPATH/src/slackibot
COPY . .
RUN go install -v

FROM alpine
RUN apk update && apk add --no-cache ca-certificates
EXPOSE 3000
COPY --from=builder /go/bin/slackibot /
CMD ["/slackibot"]