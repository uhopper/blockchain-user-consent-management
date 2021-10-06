package main

import (
	"log"
	"os"
	"path/filepath"
	"tsundoku_blockchain_ws/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "false")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		log.Println("Does")
		err = utils.PopulateWallet(wallet, "../../crypto-config/peerOrganizations/org.u-hopper.com/users/User1@org.u-hopper.com/msp/", "User1@org.u-hopper.com-cert.pem")
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	connectionFile := os.Getenv("CONNECTION_FILE")

	if connectionFile == "" {
		log.Println("using local file")
		connectionFile = "local-connection-org.yml"
	}

	ccpPath := filepath.Join(connectionFile)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("tsundokchannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("tsundokuconsent")
	result, err := contract.EvaluateTransaction("ReadConsent", "user11")

	if err != nil {
		log.Print(err)
	}
	log.Println(result)

	// log.Println("--> Create a consent")
	// result, err := contract.SubmitTransaction("UpdateConsent", "user11", "true", "asd")
	// if err != nil {
	// 	log.Fatalf("Failed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))

}
