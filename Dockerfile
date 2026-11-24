FROM golang:1.22.12-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/courses ./cmd/server

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/courses /usr/local/bin/courses
USER app
ENV HTTP_ADDR=:8080 GRPC_ADDR=:9090
EXPOSE 8080 9090
ENTRYPOINT ["/usr/local/bin/courses"]
