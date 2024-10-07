FROM golang:1.23.0 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /shareless cmd/shareless/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /shareless /shareless
COPY --from=build /app/web /web

EXPOSE ${PORT:-3000}

USER nonroot:nonroot

ENTRYPOINT ["/shareless"]
