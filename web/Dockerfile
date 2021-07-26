# Start by building the application.
FROM golang:1.16-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -a -o /go/bin/app main.go

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
ADD templates /templates
CMD ["/app"]