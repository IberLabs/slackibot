FROM golang:1.9.2

LABEL maintainer="Manuel Boira <manuelbcd@gmail.com>"

WORKDIR $GOPATH/src/slackibot

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["go-docker"]