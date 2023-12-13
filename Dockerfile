FROM golang:1.21.4 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV GOARCH=amd64

RUN go build \
    -o /go/bin/app

FROM gcr.io/distroless/base-debian11

COPY --from=build /go/bin/app /app

EXPOSE 8080

USER nonroot:nonroot

CMD ["/app"]