FROM golang:1.14-alpine AS build

WORKDIR /build

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .

RUN go build -o seiteki cmd/seiteki/main.go


FROM alpine:latest AS final

COPY --from=build /build/seiteki /bin/seiteki

RUN chmod +x /bin/seiteki

ENTRYPOINT ["/bin/seiteki"]