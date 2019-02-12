# couchdb-api

The CouchDB API runs over the Couch Go Lang dependency to connect to the CouchDB and lists all the documents in the given database.

## Input:
The CouchDB API accepts the following as environment variable.

ENV Name | Description
-- | --
LISTEN_PORT | The listening port for the CouchAPI server. This is mandatory.
COUCHDB_PORT | The DB Port number for CouchDB Host. Not mandatory if connecting via `k8s ingress`. The port number is mandatory if using the `service url`
COUCH_HOST | The database host where the CouchDB is running. A `k8s ingress` or `svc url`. This is mandatory.
SERVE_DATABASE | The database in the above CouchDB to be read served over the API. This is mandatory.

The server will serve all the data over the `/data` API. The logs are formatted using the `logrus`.

