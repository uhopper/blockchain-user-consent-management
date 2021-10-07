# Tsundoku Blockchain WS

## cUrl examples

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"ID":"user13@asds","consent": true, "privacyPolicyHash": "asdasdas"}' \
  http://localhost:5000/consent

  curl --header "Content-Type: application/json" \
  --request POST -H "apikey: 12345" \
  --data '{"ID":"user13@asds","consent": true, "privacyPolicyHash": "asdasdas"}' \
  http://localhost:5000/consent

    curl --header "Content-Type: application/json" \
  --request POST -H "apikey: 1234" \
  --data '{"ID":"user13@asds","consent": false, "privacyPolicyHash": "asdasdas"}' \
  http://localhost:5000/consent
```