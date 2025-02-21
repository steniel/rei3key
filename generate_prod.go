package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

// Simplified License struct, make sure it matches your r3/types.License
type License struct {
	LicenseId     string   `json:"licenseId"`
	ClientId      string   `json:"clientId"`
	Extensions    []string `json:"extensions"`
	LoginCount    int64    `json:"loginCount"`
	RegisteredFor string   `json:"registeredFor"`
	ValidUntil    int64    `json:"validUntil"`
}

type LicenseFile struct {
	License   License `json:"license"`
	Signature string  `json:"signature"`
}

type LicenseParams struct {
	LicenseId     string    `json:"licenseId"`
	ClientId      string    `json:"clientId"`
	Extensions    []string  `json:"extensions"`
	LoginCount    int64     `json:"loginCount"`
	RegisteredFor string    `json:"registeredFor"`
	ValidUntil    time.Time `json:"validUntil"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <license_params_file>")
		return
	}

	paramsFilename := os.Args[1]

	paramsFile, err := os.ReadFile(paramsFilename)
	if err != nil {
		panic(err)
	}

	var params LicenseParams
	err = json.Unmarshal(paramsFile, &params)
	if err != nil {
		panic(err)
	}

	validUntil := params.ValidUntil.Unix()

	// Load your private key
	privateKeyBytes, err := os.ReadFile("private.pem")
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		panic("failed to decode private key")
	}

	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// If the key is encrypted, it might fail here.
		// Try to decrypt it with an empty passphrase.
		privateKeyInterface, _, err = pkcs12.Decode(block.Bytes, "")
		if err != nil {
			panic(err)
		}
		privateKeyBytes, ok := privateKeyInterface.([]byte)
		if !ok {
			panic("failed to assert private key interface to []byte")
		}
		privateKeyInterface, err = x509.ParsePKCS8PrivateKey(privateKeyBytes)
		if err != nil {
			panic(err)
		}

	}

	privateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		panic("private key is not an RSA key")
	}

	// Create the License Data
	license := License{
		LicenseId:     params.LicenseId,
		ClientId:      params.ClientId,
		Extensions:    params.Extensions,
		LoginCount:    params.LoginCount,
		RegisteredFor: params.RegisteredFor,
		ValidUntil:    validUntil,
	}

	// Marshal the License Data to JSON
	licenseJSON, err := json.Marshal(license)
	if err != nil {
		panic(err)
	}

	// Hash the JSON Data
	hashed := sha256.Sum256(licenseJSON)

	// Sign the Hash with the Private Key
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}

	// Base64 Encode the Signature
	signatureBase64 := base64.URLEncoding.EncodeToString(signature)

	// Create the LicenseFile Structure (Matching your r3/types.LicenseFile)
	licenseFile := struct {
		License   License `json:"license"`
		Signature string  `json:"signature"`
	}{
		License:   license,
		Signature: signatureBase64,
	}

	// Marshal the complete license file into json.
	licenseFileJSON, err := json.Marshal(licenseFile)
	if err != nil {
		panic(err)
	}

	// Output the Complete License File JSON
	fmt.Println(string(licenseFileJSON))
}
