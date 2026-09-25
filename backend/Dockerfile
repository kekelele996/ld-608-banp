FROM golang:1.22-alpine AS build
WORKDIR /app

# Cache module downloads first. The default build is CGO-free (MySQL only);
# the SQLite driver is build-tagged and never linked into this image.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache wget ca-certificates
COPY --from=build /app/server ./server

ENV SERVER_PORT=3000
ENV DB_DRIVER=mysql
ENV DB_HOST=db
ENV DB_PORT=3306
EXPOSE 3000
CMD ["./server"]
