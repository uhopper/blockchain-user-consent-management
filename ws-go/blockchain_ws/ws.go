package main

import (
	"blockchain_ws/model"
	"blockchain_ws/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type BlockChainHandler struct {
	contract         *gateway.Contract
	authorizedApikey string
}

func (handler *BlockChainHandler) checkAuthentication(r *http.Request) error {
	apikey := r.Header.Get("apikey")
	if apikey == "" {
		return fmt.Errorf("Missing ApiKey header")
	}

	if apikey != handler.authorizedApikey {
		return fmt.Errorf("Not Authorized")
	}

	return nil
}

func (handler *BlockChainHandler) getConsent(w http.ResponseWriter, r *http.Request) {

	authErr := handler.checkAuthentication(r)

	if authErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%v\n", authErr)
		return
	}

	params := mux.Vars(r)
	userId := params["userId"]
	hashedUserId := utils.HashString(userId)
	log.Printf("Retrieving consent for user [%v][%v]", userId, hashedUserId)

	result, err := handler.contract.EvaluateTransaction("ReadConsent", hashedUserId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Consent for user not found\n")
		return
	}

	var consent model.Consent
	err = json.Unmarshal(result, &consent)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Unable to retrieve the content\n")
		return
	}

	consent.ID = userId

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consent)
}

func (handler *BlockChainHandler) updateConsent(w http.ResponseWriter, r *http.Request) {

	authErr := handler.checkAuthentication(r)

	if authErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%v\n", authErr)
		return
	}

	var consentRequest model.ConsentWritable

	err := json.NewDecoder(r.Body).Decode(&consentRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body is not a valid json\n")
		return
	}

	log.Printf("Update consent for user [%v]", consentRequest.ID)

	_, err = handler.contract.SubmitTransaction("UpdateConsent", utils.HashString(consentRequest.ID), strconv.FormatBool(consentRequest.UserConsent), consentRequest.PrivacyPolicyHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to Submit transaction: %v\n", err)
		log.Printf("Failed to Submit transaction: %v", err)
		return
	}

	log.Printf("Updated consent for user [%v]", consentRequest.ID)

	w.WriteHeader(http.StatusAccepted)
}

func main() {

	authorizedApikey := os.Getenv("AUTHORIZED_APIKEY")
	mspDirectory := os.Getenv("MSP_DIRECTORY")
	connectionConfigFile := os.Getenv("CONNECTION_CONFIG_FILE")

	if mspDirectory == "" {
		mspDirectory = "../../crypto-config/peerOrganizations/org.u-hopper.com/users/User1@org.u-hopper.com/msp/"
	}

	if connectionConfigFile == "" {
		log.Println("using default connection file")
		connectionConfigFile = "ws-connection-config.yml"
	}

	if authorizedApikey == "" {
		log.Fatalln("Missing AUTHORIZED_APIKEY env variable")
	}

	contract, err := utils.GetContract(mspDirectory, connectionConfigFile)

	if err != nil {
		log.Fatalf("Unable to connecto to the blockchain: %v", err)
	}

	handler := &BlockChainHandler{contract: contract, authorizedApikey: authorizedApikey}
	router := mux.NewRouter()
	router.HandleFunc("/consent/{userId}", handler.getConsent).Methods("GET")
	router.HandleFunc("/consent", handler.updateConsent).Methods("POST")
	http.Handle("/", router)
	log.Println("Server ready on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
