FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /meshgate ./cmd/meshgate

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /meshgate /meshgate
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/meshgate"]
