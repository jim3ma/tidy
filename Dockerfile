FROM alpine:3.3

MAINTAINER Jim Mar <majinjing3@gmail.com>
ENV CREATE_DATE 2016-05-03

ENV TIDY_DIR /usr/share/tidy/

ADD build/bin ${TIDY_DIR}

WORKDIR /usr/share/tidy/

EXPOSE 8089

CMD ["/usr/share/tidy/tidy"]

