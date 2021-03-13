FROM golang:1.16

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS=-mod=vendor

RUN apt-get update -y && apt-get install -y --no-install-recommends apt-utils

WORKDIR /tmp/parser

ADD . .

RUN apt-get install make git zlib1g-dev libssl-dev gperf php-cli cmake g++ -y
RUN git clone https://github.com/tdlib/td.git
RUN mkdir td/build
RUN cd td/build && cmake -DCMAKE_BUILD_TYPE=Release -DOPENSSL_ROOT_DIR=/usr/lib/ssl ..
RUN cd td/build && cmake --build . --target install
RUN go build

EXPOSE 8080

CMD ./telegram-parser -mod=vendor

