# ยังหาวิธีใช้กับ .env อยู่
FROM golang:1.19-alpine as build-base

WORKDIR /app
COPY . .
RUN go mod download
# RUN CGO_ENABLED=0 go test -v ./
RUN go build -o ./out/go-app .

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-app /app/go-app

CMD ["/app/go-app"]