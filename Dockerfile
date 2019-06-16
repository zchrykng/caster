FROM golang

ADD . /go/src/github.com/zchrykng/caster

RUN go get github.com/zchrykng/caster
RUN go install github.com/zchrykng/caster

ENTRYPOINT ["/go/bin/caster", "-root=/casts", "-userfile=/users.json"]
