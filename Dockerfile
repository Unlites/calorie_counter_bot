FROM golang:1.16.14-alpine3.14 as builder

COPY . /github.com/Unlites/calorie_counter_bot

WORKDIR /github.com/Unlites/calorie_counter_bot

RUN go mod download

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/Unlites/calorie_counter_bot/.bin/calloriecounter .
COPY --from=0 /github.com/Unlites/calorie_counter_bot/configs/ configs/

CMD ["./calloriecounter"]