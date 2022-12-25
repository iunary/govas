FROM golang as build
WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download

ADD public/ /app/public
COPY main.go /app

RUN cd /app && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o govas

FROM alpine
COPY --from=build /app/govas /app/govas
COPY --from=build /app/public /app/public
WORKDIR /app
EXPOSE 4444
CMD ["./govas"]