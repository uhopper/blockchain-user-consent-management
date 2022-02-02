package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// GetContract connect to blockchain and get an handler for the contract
func GetContract(mspDirectory string, connectionConfigFile string) (*gateway.Contract, error) {

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "false")
	if err != nil {
		log.Printf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
		return nil, err
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Printf("Failed to create wallet: %v", err)
		return nil, err
	}

	if !wallet.Exists("appUser") {
		log.Println("Does")
		err = PopulateWallet(wallet, mspDirectory, "User1@org.u-hopper.com-cert.pem")
		if err != nil {
			log.Printf("Failed to populate wallet contents: %v", err)
			return nil, err
		}
	}

	ccpPath := filepath.Join(connectionConfigFile)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Printf("Failed to connect to gateway: %v", err)
		return nil, err
	}
	defer gw.Close()

	network, err := gw.GetNetwork("consentmanagementchannel")
	if err != nil {
		log.Printf("Failed to get network: %v", err)
		return nil, err
	}

	contract := network.GetContract("consentmanagementcontract")
	return contract, nil
}
