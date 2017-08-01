FROM golang:1.8

MAINTAINER Timothy S. Williams <tiwillia@redhat.com>

USER nobody

RUN mkdir -p /go/src/github.com/tiwillia/gomr
WORKDIR /go/src/github.com/tiwillia/gomr

COPY . /go/src/github.com/tiwillia/gomr
RUN go-wrapper download github.com/tiwillia/gomr/cmd/gomr && go-wrapper install github.com/tiwillia/gomr/cmd/gomr

CMD ["go-wrapper", "run", "-logtostderr"]
