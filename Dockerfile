FROM alpine:3.17 AS build

RUN apk update
RUN apk upgrade
RUN apk add --update go

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o /portcheck

FROM alpine:3.17

WORKDIR /

COPY --from=build /portcheck /portcheck

CMD ["/portcheck"]