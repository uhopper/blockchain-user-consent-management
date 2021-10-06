package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

//Populate the wallet
func PopulateWallet(wallet *gateway.Wallet, mspFolderPath string, user string) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(mspFolderPath)

	certPath := filepath.Join(credPath, "signcerts", user)
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("UHopperOrgMSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
