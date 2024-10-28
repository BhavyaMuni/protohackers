ARG GO_VERSION=1

FROM golang:${GO_VERSION}-bookworm as builder

ENV QUESTION_NUMBER=0
WORKDIR /usr/src/app
COPY $QUESTION_NUMBER/go.mod ./
RUN go mod download && go mod verify
COPY $QUESTION_NUMBER/. .
RUN go build -v -o /run-app .


FROM debian:bookworm

COPY --from=builder /run-app /usr/local/bin/
EXPOSE  10000
CMD ["run-app"]
