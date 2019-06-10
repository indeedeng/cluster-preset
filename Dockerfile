FROM golang:1.12

RUN go get golang.org/x/lint/golint

ENV GO111MODULE on

WORKDIR /go/src/cluster-preset

COPY . .

RUN make deps test install

ENTRYPOINT [ "cluster-preset" ]
