FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN GOTOOLCHAIN=local go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/monitor ./cmd/monitor

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/monitor /monitor
EXPOSE 8080
ENTRYPOINT ["/monitor"]
