#go test ./... -covermode=count --coverprofile='coverage.out'

./certgen.sh
echo "..."
docker build -t couchdbdrv -f ./docker/Dockerfile ./docker