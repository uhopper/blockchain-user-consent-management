# Chaincode (GoLang)

This folder contains the implementation of the chaincode (SmartContract) used in the *user-consent-management* blockchain.

This chaincode define two transactions:

- `ReadConsent`: which allows to check the consent for a given user
- `UpdateConsent`: which allows to update the consent for the given user

The transaction definition is located inside the `smartcontract.go` file in the `chaincode` directory.