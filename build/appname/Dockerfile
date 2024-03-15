FROM golang:1.22 AS build

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download

COPY ./pkg ./pkg
COPY ./cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server /go/src/app/cmd/appname


FROM gcr.io/distroless/static-debian12
# FROM alpine:latest # use it for development purposes to provide shell access

COPY --from=build /go/src/app/server /

EXPOSE 1312

CMD ["/server"]
