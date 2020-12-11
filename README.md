# couchdb
Driver for Apache CouchDB&reg; 3.x
Supports ssl connection and cookie authentication.

make build: creates certificates and builds docker image.
make start: starts CouchDB database in docker container, exposes standard couch ports: 5984 and 6984 on 5300 and 6300 respectively.
make tests: runs tests, generates code coverage.
make stop: removes docker app.