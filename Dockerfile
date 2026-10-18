FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/paypilot ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/paypilot /paypilot
EXPOSE 8080
ENTRYPOINT ["/paypilot"]
