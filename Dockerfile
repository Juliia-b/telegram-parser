#  docker build -t tg-parser .

FROM golang:1.16

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS=-mod=vendor

RUN apt-get update -y && apt-get install -y --no-install-recommends apt-utils

WORKDIR /tmp/parser

ADD . .

RUN apt-get install make
RUN apt-get install git
RUN apt-get install -y zlib1g-dev
RUN apt-get install -y libssl-dev
RUN apt-get install gperf
RUN apt-get install -y php-cli
RUN apt-get install -y cmake
RUN apt-get install -y g++
RUN git clone https://github.com/tdlib/td.git
RUN mkdir td/build
RUN cd td/build && cmake -DCMAKE_BUILD_TYPE=Release -DOPENSSL_ROOT_DIR=/usr/lib/ssl ..
RUN cd td/build && cmake --build . --target install
RUN go build

EXPOSE 8080

CMD ./telegram-parser -mod=vendor

