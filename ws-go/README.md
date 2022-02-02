# Blockchain WS

This folder contains the implementation of the WS which allows easy access to the *user-consent-management* blockchain. This web service will handle the authentication and the connection with the blockchain and expose two endpoints for handling the data in the blockchain.

The OpenApi documentation of the endpoints is availble in the [documentation/openapi.yaml](documentation/openapi.yaml) file

## Setup

### Environmental Variables

- *AUTHORIZED_APIKEY*: The apikey used for ws user authentication.
- *MSP_DIRECTORY*:  path to the `msp` folder of the blockchain user dedicated to the webservice
- *CONNECTION_CONFIG_FILE*: The path to the `ws-connection-config.yml` configuration file.

## Run the server

First export the required environmental variables:

```bash
export AUTHORIZED_APIKEY=<your apikey>
export MSP_DIRECTORY=/path/to/msp/folder
export CONNECTION_CONFIG_FILE=/path/to/config/file
```

Run the webserice:

```bash
cd blockchain_ws
go run ws.go
```

The web service will starts on port `5000`

### cUrl examples

```bash
  curl --header "Content-Type: application/json" \
  --request POST -H "apikey: your apikey" \
  --data '{"ID":"user@example.com","consent": true, "privacyPolicyHash": "AjZ-EVOt5xRCywnjONkgtHyDi71etWE8DV0byVOBEjw"}' \
  http://localhost:5000/consent
```

```bash
    curl --header "Content-Type: application/json" \
  --request POST -H "apikey: your apikey" \
  --data '{"ID":"user@example.com","consent": false, "privacyPolicyHash": "AjZ-EVOt5xRCywnjONkgtHyDi71etWE8DV0byVOBEjw"}' \
  http://localhost:5000/consent
```

```bash
  curl --header "Content-Type: application/json"  -H "apikey: your apikey" http://localhost:5000/consent/stefano.tavonatti+test@u-hopper.com
```
