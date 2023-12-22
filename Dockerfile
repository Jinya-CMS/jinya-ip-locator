FROM library/golang:1.21-alpine as build
WORKDIR /app
COPY . .

RUN go build -o /jinya-ip-locator

FROM library/alpine:latest

COPY --from=build /jinya-ip-locator /jinya-ip-locator

CMD ["/jinya-ip-locator"]
