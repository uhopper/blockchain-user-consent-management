package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	model "mortadella-market/model"
	stateutils "mortadella-market/stateutils"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//SmartContract the smart contract struct
type SmartContract struct {
}

var logger = shim.NewLogger("tsundokuconsent")

//Init method is called when the Smart Contract "tsundokuconsent" is instantiated by the blockchain network Best practice is to have any Ledger initialization in separate function -- see initLedger()
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//Invoke method is called as a result of an application request to run the Smart Contract "tsundokuconsent". The calling application program has also specified the particular smart contract function to be called, with arguments
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "tsundokuconsent:getCompany" {
		return s.getCompany(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "tsundokuconsent:createComapany" {
		return s.createComapany(APIstub, args)
	} else if function == "tsundokuconsent:listCompanies" {
		return s.listCompanies(APIstub, args)
	} else if function == "tsundokuconsent:trade" {
		return s.trade(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name. [" + function + "]")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("Grifo contract initialized")
	return shim.Success(nil)
}

func (s *SmartContract) getCompany(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	company, err := stateutils.GetCompanyFromWorldState(APIstub, args[0])

	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve the company. Error: %s", err))
	}

	companyDTO := model.ComapanyDTO{}
	companyDTO.FromCompany(*company)

	logger.Infof("Retrieved comapany [%s]", company)

	companyAsBytes, _ := json.Marshal(companyDTO)

	return shim.Success(companyAsBytes)
}

func (s *SmartContract) createComapany(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	identity, errIdentity := APIstub.GetCreator()

	if errIdentity != nil {
		logger.Errorf("Unable to retrieve the identity [%v]", errIdentity)
		return shim.Error("Unable to retrieve the identity ")
	}

	mortadellaUnit, errParseint := strconv.ParseInt(args[2], 10, 64)

	if errParseint != nil {
		logger.Errorf("Unable to parse the string [%s] to int", args[2])
		return shim.Error(fmt.Sprintf("Unable to parse the string [%s] to int", args[2]))
	}

	var company = model.Company{State: model.State{ID: args[0]}, Name: args[1], Identity: base64.URLEncoding.EncodeToString(identity), MortadellaUnits: mortadellaUnit}

	err := stateutils.SaveState(APIstub, company)

	if err != nil {
		logger.Errorf("Unabel to save the company [%v]", company)
		return shim.Error(fmt.Sprintf("Unable to save the company [%v]", company))
	}

	companyDTO := model.ComapanyDTO{}
	companyDTO.FromCompany(company)

	companyAsBytes, errMarshaling := json.Marshal(companyDTO)

	if errMarshaling != nil {
		logger.Errorf("Unable to marshal the company [%v]", company)
		return shim.Error("Unable to create the company")
	}

	logger.Infof("Saved company [%v]", company)
	return shim.Success(companyAsBytes)
}

func (s *SmartContract) trade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	idSource := args[0]
	idDestination := args[1]

	quantity, errParseint := strconv.ParseInt(args[2], 10, 64)

	if errParseint != nil {
		logger.Errorf("Unable to parse the string [%s] to int", args[2])
		return shim.Error(fmt.Sprintf("Unable to parse the string [%s] to int", args[2]))
	}

	identity, errIdentity := APIstub.GetCreator()

	if errIdentity != nil {
		logger.Errorf("Unable to retrieve the identity [%v]", errIdentity)
		return shim.Error("Unable to retrieve the identity ")
	}

	companySource, errCompany := stateutils.GetCompanyFromWorldState(APIstub, idSource)

	if errCompany != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve the company. Error: %s", errCompany))
	}

	companyDestination, errDestination := stateutils.GetCompanyFromWorldState(APIstub, idDestination)

	if errDestination != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve the company. Error: %s", errDestination))
	}

	//check if the company which try to transfer the mortadella is the owner of the account

	stringIdentity := base64.URLEncoding.EncodeToString(identity)

	if companySource.Identity != stringIdentity {
		logger.Errorf("[%s] not equal to [%s]", companySource.Identity, string(identity))
		return shim.Error("Not Allowed")
	}

	errTransfer := transferMortadella(APIstub, *companySource, *companyDestination, quantity)

	if errTransfer != nil {
		return shim.Error(fmt.Sprintf("Unable to tranfer the mortadella, error: %s", errTransfer))
	}

	return shim.Success([]byte("OK"))
}

func transferMortadella(APIstub shim.ChaincodeStubInterface, source model.Company, destination model.Company, quantity int64) error {

	//check if source company has enough mortadella

	if source.MortadellaUnits >= quantity {
		source.MortadellaUnits -= quantity
		destination.MortadellaUnits += quantity

		errSource := stateutils.SaveState(APIstub, source)

		if errSource != nil {
			return fmt.Errorf("Unable to update the source company")
		}

		errDestination := stateutils.SaveState(APIstub, destination)

		if errDestination != nil {
			return fmt.Errorf("Unable to update the destination company")
		}

		logger.Infof("Sucessfully transfer mortadella [%v] units from [%s] to [%s]", quantity, source.ID, destination.ID)

		return nil

	}

	return fmt.Errorf("Not enough mortadella")
}

func (s *SmartContract) listCompanies(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	companyList, err := stateutils.ListCompaniesFromWorldState(APIstub)

	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to retrieve the company list. Error: %s", err))
	}

	companyDTOList := []model.ComapanyDTO{}

	for _, c := range companyList {
		companyDTO := model.ComapanyDTO{}
		companyDTO.FromCompany(c)
		companyDTOList = append(companyDTOList, companyDTO)
	}

	companyListAsBytes, errMarshal := json.Marshal(companyDTOList)

	if errMarshal != nil {
		logger.Errorf("Unable to marshal the company list. Error: %s", errMarshal)
		return shim.Error(fmt.Sprintf("Unable to retrieve the company list"))
	}

	logger.Infof("Retrieved [%d] companies", len(companyList))
	return shim.Success(companyListAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))

	logger.SetLevel(shim.LogDebug)
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
