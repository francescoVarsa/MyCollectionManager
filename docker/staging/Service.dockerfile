FROM alpine:latest

WORKDIR /api

COPY ./../backend/API/build/services ./services
COPY ./../backend/API/.env ./.env

RUN chmod +x ./services

EXPOSE 8080

CMD ["./services", "-envFile", "./.env"]