FROM golang:1.20.2-alpine3.17 AS build

WORKDIR /usr/src/portcheck

COPY go.mod ./
COPY internal ./internal
COPY cmd ./cmd

RUN go build -o /usr/local/bin/portcheck cmd/portcheck/main.go

FROM alpine:3.17

COPY --from=build /usr/local/bin/portcheck /portcheck

CMD ["/portcheck"]