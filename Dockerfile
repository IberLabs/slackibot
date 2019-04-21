FROM golang:1.11.9-alpine3.9 AS builder
LABEL maintainer="Manuel Boira <manuelbcd@gmail.com>"
RUN apk update
RUN apk upgrade
RUN apk add --no-cache git gcc g++ musl-dev bash
WORKDIR /go/src/slackibot
COPY . .
RUN go get -d -v /go/src/slackibot/cmd/slackibot/
#RUN go install -v /go/src/slackibot/cmd/slackibot/
RUN /bin/bash /go/src/slackibot/scripts/build.sh

FROM alpine:3.9
LABEL maintainer="Manuel Boira <manuelbcd@gmail.com>"
RUN apk update && apk add --no-cache ca-certificates libc6-compat
EXPOSE 3000
RUN mkdir /slackibot
COPY --from=builder /go/src/slackibot/bin /slackibot/
RUN chmod +x /slackibot/slackibot
RUN chmod -R +wxr /slackibot
#ENTRYPOINT ["/slackibot"]
CMD ["./slackibot"]