FROM golang:1.16-alpine AS builder
ARG user=test
WORKDIR /app
COPY . .
ENV GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn
RUN go build -o reservation_thxx_go \
    && go build -o reservation_thxx_go_external ./external

FROM ubuntu:20.04 AS runner
RUN apt-get update && apt-get -y install tar cron mongo-tools
WORKDIR /app
COPY --from=builder /app/ .
EXPOSE 9000
CMD crontab tools/crontab.job \
    && ./reservation_thxx_go --staging=false --web=:9000