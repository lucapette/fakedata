FROM alpine

RUN apk add --no-cache strace

ADD ./fakedata .
