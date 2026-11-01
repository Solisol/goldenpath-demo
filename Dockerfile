# goldenpath-api: the golden-path service, with the Markdown articles built in.
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY articles.go ./
COPY articles/ ./articles/
COPY cmd/goldenpath-api/ ./cmd/goldenpath-api/
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /goldenpath-api ./cmd/goldenpath-api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /goldenpath-api /goldenpath-api
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/goldenpath-api"]
