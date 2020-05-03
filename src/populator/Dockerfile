FROM golang:latest as builder

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o populator .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
COPY --from=builder /app/populator /app/
WORKDIR /app/
COPY schema.sql /app/
COPY indexes.sql /app/
