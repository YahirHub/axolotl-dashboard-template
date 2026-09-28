FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/axolotl .

FROM alpine:3.22

RUN addgroup -S axolotl && adduser -S -G axolotl axolotl

WORKDIR /app

COPY --from=builder /out/axolotl /app/axolotl
COPY --from=builder /src/public /app/public
COPY --from=builder /src/resources /app/resources

RUN mkdir -p /app/storage/logs && chown -R axolotl:axolotl /app

USER axolotl

EXPOSE 3000

ENTRYPOINT ["/app/axolotl"]
