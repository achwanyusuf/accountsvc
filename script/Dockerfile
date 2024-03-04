FROM golang:1.20
WORKDIR /app
COPY ./conf/conf.yaml conf/
COPY ./migrations migrations/
COPY ./build/app .

EXPOSE 8081

# Run
CMD ["./app"]