FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
ENV GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn
RUN go build -o reservation_thxx_go \
    && go build -o reservation_thxx_go_external ./external

FROM alpine AS runner
WORKDIR /app
COPY --from=builder /app/ .
EXPOSE 9000
CMD [ "./reservation_thxx_go" ]