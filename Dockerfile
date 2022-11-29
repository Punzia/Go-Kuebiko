FROM golang:1.19.3 AS build

LABEL authors="Joseph Alvarenga and Rapunzel"
LABEL version='0.0.1'

WORKDIR  /src

COPY . .
RUN go mod init polling-service
RUN go get .
RUN GOOS=linux CGO_ENABLED=0 go build -o kuebiko

# Deploy stage
FROM alpine:3.15

WORKDIR /usr/bin

COPY --from=build /src/ .

EXPOSE 8080

ENTRYPOINT ["go-kuebiko"]