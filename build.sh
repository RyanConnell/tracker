export GOARCH=amd64
export GOOS=linux

set -e

echo "Building tracker-frontend"
go build -o bin/tracker-frontend cmd/frontend/frontend.go
echo "Building tracker-backend"
go build -o bin/tracker-backend cmd/backend/backend.go
echo "Building tracker-scraper"
go build -o bin/tracker-scraper cmd/scraper/scraper.go

chmod a+x bin/*

echo "Building docker image"
docker build . -t tracker
