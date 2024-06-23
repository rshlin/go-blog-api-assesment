FROM golang:1.19 AS build
WORKDIR /go/src
COPY api ./api
COPY out/go/main.go .
COPY out/go/go.sum .
COPY out/go/go.mod .

ENV CGO_ENABLED=0

RUN go build -o app .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/app ./
EXPOSE 8080/tcp
ENTRYPOINT ["./app technical --address 0.0.0.0 --port 8080"]
