FROM golang

ADD ./users.json /users.json

RUN go get github.com/zchrykng/gocaster/cmd/caster
RUN go install github.com/zchrykng/gocaster/cmd/caster

ENTRYPOINT ["/go/bin/caster", "-root=/casts", "-userfile=/users.json"]
