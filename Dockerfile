FROM golang:1.18.2-alpine as build
WORKDIR /src/app
COPY . .
RUN go mod download
RUN go build -o /bin/go-api

FROM alpine
COPY --from=build /bin/go-api /bin/go-api
ENTRYPOINT ["/bin/go-api"]
