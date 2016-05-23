FROM golang:1.6-alpine

MAINTAINER Jim Mar <majinjing3@gmail.com>
ENV CREATE_DATE 2016-05-03

ENV TIDY_DIR /usr/share/tidy/
ENV BUILD_DIR /go/src/github.com/jim3mar/tidy/

ADD . ${BUILD_DIR}

WORKDIR /usr/share/tidy/

EXPOSE 8089

CMD ["/usr/share/tidy/tidy"]
