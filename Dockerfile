FROM golang:1.24 as builder



WORKDIR /app

COPY . .





RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /test-task cmd/main.go


FROM alpine:3.21


WORKDIR /app

COPY  --from=builder /test-task .

 
COPY --from=builder /app/db  ./db 



EXPOSE 8090

CMD [ "./test-task" ]