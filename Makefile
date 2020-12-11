build:
	./prepare.sh
dockerbuild:
	docker build -t couchdbdrv -f ./docker/Dockerfile ./docker
start:
	docker run -d -l couchdbdrv1  -p5300:5984 -p6300:6984 -v $(CURDIR)/docker/etc:/opt/couchdb/etc/local.d couchdbdrv
stop:
	docker rm -f $$( docker ps -qaf "label=couchdbdrv1")
tests:
	go test ./... -covermode=count --coverprofile='coverage.out'
	go tool cover -html coverage.out 