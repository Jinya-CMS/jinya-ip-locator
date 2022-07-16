FROM quay.imanuel.dev/dockerhub/library---golang:1.18-alpine as build
WORKDIR /app
COPY . .

RUN go build -o /jinya-ip-locator

FROM quay.imanuel.dev/dockerhub/library---alpine:latest

COPY --from=build /jinya-ip-locator /jinya-ip-locator

CMD ["/jinya-ip-locator"]
