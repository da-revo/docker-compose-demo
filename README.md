CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

docker-compose build
docker-compose up -d
docker-compose stop