FROM golang:alpine AS build
WORKDIR /go/src/duoc-plus
COPY . .
RUN go build -o /go/bin/duoc-plus cmd/main.go

FROM scratch
COPY --from=build etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/duoc-plus /go/bin/duoc-plus
ENTRYPOINT ["/go/bin/duoc-plus"]