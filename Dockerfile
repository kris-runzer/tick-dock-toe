FROM golang:1.7-alpine

ADD ./ /go/src/github.com/kris-runzer/tick-dock-toe

WORKDIR /go/src/github.com/kris-runzer/tick-dock-toe

RUN go install

ENTRYPOINT ["tick-dock-toe"]
