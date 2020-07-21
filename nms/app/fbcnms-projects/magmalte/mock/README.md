Mock server
- Enables us to mock the json responses from orc8r
https://github.com/typicode/json-server

1. create certs in mock directory (magmalte/mock)
openssl req -nodes -new -x509 -keyout .cache/mock_server.key -out .cache/mock_server.cert

2. Add the mock certs and API server info to the .env file i.e (magmalte/.env)
# mock
+API_HOST=https://<hostname>:3001
+API_CERT_FILENAME=<WORKING_DIR>/magmalte/mock/.cache/mock_server.cert
+API_PRIVATE_KEY_FILENAME=<WORKING_DIR>/magmalte/mock/.cache/mock_server.key

2. run the mock service
yarn run mockserver
