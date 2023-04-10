# FROM golang:1.15.8-alpine
FROM golang:1.17-alpine

RUN \
    apk add --no-cache bash git openssh && \
    apk --no-cache add curl && \
    apk --no-cache add vim && \
    apk --no-cache add procps-dev && \
    apk --no-cache add busybox-extras

ADD ./ /app
RUN cd /app
WORKDIR /app/cmd
RUN go build -o main .
CMD ["/app/cmd/main"]