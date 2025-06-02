############################
# STEP 1: Build executable
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/main .

############################
# STEP 2: Create final image
############################
FROM alpine:3.22

ENV PORT=8080
ENV GIN_MODE=release

WORKDIR /go
COPY --from=builder /go/main .
# COPY public /go/public

EXPOSE 8081
ENTRYPOINT ["./main"]