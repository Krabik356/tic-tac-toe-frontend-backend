FROM golang
WORKDIR /server
COPY . .
EXPOSE 8080
RUN go build -o app ./server/cmd/
CMD ["./app"]