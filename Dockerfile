FROM golang:1.11 AS builder
LABEL maintainer="Manuel Boira <manuelbcd@gmail.com>"
#RUN apk update
#RUN apk upgrade
#RUN apk add --no-cache git gcc g++ musl-dev bash

WORKDIR /go/src/slackibot
COPY . .
RUN go get -d -v /go/src/slackibot/cmd/slackibot/
#RUN go install -v /go/src/slackibot/cmd/slackibot/
RUN /bin/bash /go/src/slackibot/scripts/build.sh

#------------------------------------------------
# STAGE 2
#------------------------------------------------
FROM alpine:3.9
LABEL maintainer="Manuel Boira <manuelbcd@gmail.com>"
RUN apk update && apk add --no-cache ca-certificates libc6-compat
RUN mkdir /slackibot
COPY --from=builder /go/src/slackibot/bin /slackibot/

# File permissions
RUN chmod +x /slackibot/slackibot
RUN chmod -R +wxr /slackibot

# Create a user group 'slackibot'
RUN addgroup -g 1000 -S slackibot && \
    adduser -u 1000 -S slackibot -G slackibot

RUN chown -R slackibot:slackibot /slackibot

# Switch to 'slackibot' user
USER slackibot
EXPOSE 3000

#ENTRYPOINT ["/slackibot"]
WORKDIR /slackibot
CMD ["./slackibot"]