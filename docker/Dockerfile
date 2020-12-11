FROM couchdb
COPY etc/local.ini /opt/couchdb/etc/local.d
ADD etc/*.crt /opt/couchdb/etc/cert/
ADD etc/*.key /opt/couchdb/etc/cert/
ADD etc/*.csr /opt/couchdb/etc/cert/
ADD etc/*.req /opt/couchdb/etc/cert/
