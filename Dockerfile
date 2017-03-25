FROM golang:1.8
MAINTAINER Anytime <936269579@qq.com>
COPY ./src /go/src
ADD ./config.ini /go/config.ini
ADD ./config.ini /usr/bin/config.ini
RUN go get -d -v nmdim.net/opentab
RUN CGO_ENABLED=0 go install -a nmdim.net
EXPOSE 80
CMD /go/bin/app
