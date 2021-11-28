set -e

echo "Starting tracker-backend"
docker run --name tracker-backend -p 8081:8081 --network tracker-network -d tracker /bin/tracker-backend

echo "Waiting 3 seconds before starting frontend"
sleep 3

echo "Starting tracker-frontend"
docker run --name tracker-frontend -p 8080:8080 --network tracker-network -d tracker /bin/tracker-frontend
