FROM golang:1.18 as build
WORKDIR /app
COPY . /app
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go test -v ./...
RUN CGO_ENABLED=0 go build -o swagger-ui . && \
    chmod +x ./swagger-ui

FROM scratch as production
EXPOSE 80
COPY --from=build /app/swagger-ui /
ENTRYPOINT ["/swagger-ui"]