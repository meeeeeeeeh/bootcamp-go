FROM golang:1.22.2 AS build

WORKDIR /app

COPY . .
RUN go mod tidy

#RUN go build -v -o /app/main ./cmd/restaraunt_service/main.go
#хз почему но без этих флагов не работало
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/restaraunt_service/main.go

FROM alpine:latest
RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=build /app/main .
RUN chmod +x /root/main

COPY --from=build /app/internal/config ./internal/config
COPY --from=build /app/template ./template
#надо бы убрать энв файл отсюда
COPY --from=build /app/.env /root/.env


EXPOSE 8080

CMD ["./main"]

