package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Consent to be stored in the blockchain.
// ID is the hash of the user email
// UserConsent is boolean that specify if the user has given its consent
// LastUpdate is contain the epoch time (in seconds) of the last update of the asset
type Consent struct {
	ID                string `json:"ID"`
	UserConsent       bool   `json:"consent"`
	PrivacyPolicyHash string `json:"privacyPolicyHash"`
	LastUpdate        int64  `json:"lastUpdate"`
}

//UpdateConsent update the user consent
func (s *SmartContract) UpdateConsent(ctx contractapi.TransactionContextInterface, id string, consent bool, privacyPolicyHash string) error {

	now := time.Now()

	asset := Consent{
		ID:                id,
		UserConsent:       consent,
		PrivacyPolicyHash: privacyPolicyHash,
		LastUpdate:        now.Unix(),
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadConsent read the user consent
func (s *SmartContract) ReadConsent(ctx contractapi.TransactionContextInterface, id string) (*Consent, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Consent
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}
