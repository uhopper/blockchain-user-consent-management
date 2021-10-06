package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tsundoku_blockchain_ws/model"
	"tsundoku_blockchain_ws/utils"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type BlockChainHandler struct {
	contract *gateway.Contract
}

func (handler *BlockChainHandler) getConsent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	log.Printf("Retrieving consent for user [%v]", userId)

	// consent := &model.Consent{ID: "asd", UserConsent: true, PrivacyPolicyHash: "fgh", LastUpdate: 12}
	result, err := handler.contract.EvaluateTransaction("ReadConsent", userId)

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

	json.NewEncoder(w).Encode(consent)
}

func (handler *BlockChainHandler) updateConsent(w http.ResponseWriter, r *http.Request) {

	var consentRequest model.ConsentWritable

	err := json.NewDecoder(r.Body).Decode(&consentRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request body is not a valid json\n")
		return
	}

	log.Printf("Update consent for user [%v]", consentRequest.ID)

	_, err = handler.contract.SubmitTransaction("UpdateConsent", consentRequest.ID, strconv.FormatBool(consentRequest.UserConsent), consentRequest.PrivacyPolicyHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to Submit transaction: %v\n", err)
		log.Printf("Failed to Submit transaction: %v", err)
		return
	}

	log.Printf("Updated consent for user [%v]", consentRequest)

	w.WriteHeader(http.StatusAccepted)
}

func main() {

	contract, err := utils.GetContract()

	if err != nil {
		log.Fatalf("Unable to connecto to the blockchain: %v", err)
	}

	handler := &BlockChainHandler{contract}
	router := mux.NewRouter()
	router.HandleFunc("/consent/{userId}", handler.getConsent).Methods("GET")
	router.HandleFunc("/consent", handler.updateConsent).Methods("POST")
	http.Handle("/", router)
	log.Println("Server ready on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
