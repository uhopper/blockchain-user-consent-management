package stateutils

import (
	"encoding/json"
	"fmt"
	model "mortadella-market/model"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mortadellamarket.utils")

//SaveState save a new state in the blockchain
func SaveState(APIstub shim.ChaincodeStubInterface, state model.StateInterface) error {
	stateAsBytes, errMarshal := json.Marshal(state)

	if errMarshal != nil {
		logger.Errorf("Unable to marshal the state [%v]. Error: %s", state, errMarshal)
		return fmt.Errorf("Unable to marshal the state [%v]", state)
	}

	compositeKey, errCompositeKey := APIstub.CreateCompositeKey(state.GetIndex(), state.GetKey())

	if errCompositeKey != nil {
		return fmt.Errorf("Unable to create the composite key")
	}

	err := APIstub.PutState(compositeKey, stateAsBytes)

	if err != nil {
		return fmt.Errorf("Unable to save the state. Error: %s", err)
	}

	return nil
}

//GetCompanyFromWorldState allow to retrieve a company from the world state
func GetCompanyFromWorldState(APIstub shim.ChaincodeStubInterface, companyID string) (*model.Company, error) {

	compositeKey, errCompositeKey := APIstub.CreateCompositeKey(model.CompanyIndexString, model.GetCompanyKey(companyID))

	if errCompositeKey != nil {
		return nil, fmt.Errorf("Unable to create the composite key")
	}

	companyAsBytes, err := APIstub.GetState(compositeKey)

	if err != nil || len(companyAsBytes) == 0 {
		logger.Errorf("Error during company retrival retrival: %s", companyID)
		return nil, fmt.Errorf("Company " + companyID + " does not exists")
	}

	company := model.Company{}
	err1 := json.Unmarshal(companyAsBytes, &company)

	if err1 != nil {
		logger.Errorf("Unable to parse byte data [%s] in an company object", err1)
		return nil, fmt.Errorf("Unable to parse byte data in an company object")
	}

	return &company, nil
}

//ListCompaniesFromWorldState get the list of all companies
func ListCompaniesFromWorldState(APIstub shim.ChaincodeStubInterface) ([]model.Company, error) {
	resultIterator, err := APIstub.GetStateByPartialCompositeKey(model.CompanyIndexString, []string{})

	var companyList []model.Company = []model.Company{}

	if err != nil {
		logger.Errorf("Unable to retrieve the company list. Error: %s", err)
		return nil, fmt.Errorf("Unable to retrieve the company list. Error: %s", err)
	}

	for resultIterator.HasNext() {
		companyKv, errIt := resultIterator.Next()

		if errIt != nil {
			logger.Errorf("Unable to retrieve the company list. Error: %s", errIt)
			return nil, fmt.Errorf("Unable to retrieve the company list. Error: %s", errIt)
		}

		company := model.Company{}

		errMarshal := json.Unmarshal(companyKv.GetValue(), &company)

		if errMarshal != nil {
			logger.Errorf("Unable to retrieve the instituion list. Error: %s", errMarshal)
			return nil, fmt.Errorf("Unable to retrieve the instituion list. Error: %s", errMarshal)
		}

		companyList = append(companyList, company)
	}

	return companyList, nil
}
