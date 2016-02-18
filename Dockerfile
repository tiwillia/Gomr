FROM golang:1.5

MAINTAINER Timothy S. Williams <tiwillia@redhat.com>

USER nobody

RUN mkdir -p /go/src/github.com/tiwillia/gomr
WORKDIR /go/src/github.com/tiwillia/gomr

COPY . /go/src/github.com/tiwillia/gomr
RUN go-wrapper download && go-wrapper install

CMD ["go-wrapper", "run"]
