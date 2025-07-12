FROM docker.io/library/golang:1.24-alpine as build
WORKDIR /app
COPY . .

RUN go build -o /jinya-ip-locator

FROM docker.io/library/alpine:latest

COPY --from=build /jinya-ip-locator /jinya-ip-locator

CMD ["/jinya-ip-locator"]
