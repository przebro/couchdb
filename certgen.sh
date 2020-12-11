#!/bin/bash
#generate keys for rootca server and client
openssl genrsa -out ./docker/etc/root_ca.key 4096
openssl genrsa -out ./docker/etc/couchdrv.key 2048
openssl genrsa -out ./docker/etc/client.key 2048

#create root certficiate
openssl req -x509 -new -nodes -key ./docker/etc/root_ca.key -sha256 -days 1024 \
-out ./docker/etc/root_ca.crt \
-subj "/C=US/ST=State/L=City/O=Couchdbdrv/OU=DEV/CN=localuser/emailAddress=localuser@localhost.com"
#create signing request
openssl req -new -key ./docker/etc/couchdrv.key -out ./docker/etc/couchdrv.csr \
-subj "/C=US/ST=State/L=City/O=Couchdbsrv/OU=DEV/CN=localuser/emailAddress=localuser@localhost.com"
#create certificate for server
openssl x509 -req -in ./docker/etc/couchdrv.csr -CA ./docker/etc/root_ca.crt -CAkey ./docker/etc/root_ca.key \
-CAcreateserial -out ./docker/etc/couchdrv.crt -days 1024 -sha256 \

#create certificate for client
openssl req -new -key ./docker/etc/client.key -out ./docker/etc/client.req \
-subj "/C=US/ST=State/L=City/O=couchdbcli/OU=DEV/CN=localuser/emailAddress=localuser@localhost.com"
openssl x509 -req -in ./docker/etc/client.req -CA ./docker/etc/root_ca.crt -CAkey ./docker/etc/root_ca.key \
-set_serial 101010 -extensions client -days 1024 -out ./docker/etc/client.crt

#Subject: C = PL, ST = State, L = Lodz, O = Couchdrv, OU = section, CN = localuser, emailAddress = localuser@localhost.com
#Subject: C = US, ST = State, L = City, O = CouchDB driver ltd., OU = DEV, CN = localuser, emailAddress = localuser@localhost.com

#Subject: emailAddress = user@host.com, C = US, ST = State, L = City, O = CouchDB driver ltd., OU = DEV, CN = localhost
#Subject: C = PL, ST = State, L = City, O = Company ltd, OU = section, CN = localuser, emailAddress = user@localhost.com