FROM golang:1.24-bookworm AS build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v ./...
RUN go test -v ./...

WORKDIR /go/src/app/examples
RUN CGO_ENABLED=1 go build -o=/go/bin/app ./mp3-to-wav

FROM gcr.io/distroless/base-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]
