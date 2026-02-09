FROM golang:alpine

WORKDIR /app

RUN apk add --no-cache git bash

COPY go.mod go.sum ./

# Force go version to 1.24 to match the builder image and bypass host-side auto-upgrades
RUN sed -i 's/go 1.2[5-9]/go 1.24/' go.mod

# ติดตั้ง air + swag
RUN go mod download && \
    go install github.com/air-verse/air@v1.52.2 && \
    go install github.com/swaggo/swag/cmd/swag@v1.8.12

# copy source (ถ้าใช้ volume mount ก็จะทับอีกทีตอน runtime)
COPY . .

RUN mkdir -p docs tmp

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
