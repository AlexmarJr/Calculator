FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go test ./... && go vet ./...
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /server ./cmd/api

FROM scratch AS runtime
COPY --from=build /server /server
USER 65532:65532
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/server"]
