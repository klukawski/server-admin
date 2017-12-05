package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/microservice"
	"crypto/rsa"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
	"encoding/asn1"
	"io/ioutil"
)

func handleTest(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Hello World!")
}

func loadPrivate(filename string) *rsa.PrivateKey {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error while reading local key file: ", err)
	}
	privateKey, _ := pem.Decode(bytes)

	key, err := x509.ParsePKCS1PrivateKey(privateKey.Bytes)
	if err != nil {
		log.Fatal("Error, local private key is invalid, delete ./keys/local and retry to regenerate: ", err)
	}

	return key
}

func loadExternal(filename string) *rsa.PublicKey {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error while reading external key file: ", err)
	}
	public, _ := pem.Decode(bytes)

	key, err := x509.ParseCertificate(public.Bytes)
	if err != nil {
		log.Fatal("Error, external key is invalid, replace ./keys/external.pub with a valid certificate: ", err)
	}

	return key.PublicKey.(*rsa.PublicKey)
}

func generateAndSaveKeypair() {
	random := rand.Reader
	bitsize := 2048

	key, err := rsa.GenerateKey(random, bitsize)
	if err!=nil {
		log.Fatal("Fatal error while generating a keypair: ", err)
	}

	privateOutFile, err := os.Create("keys/local")
	if err != nil {
		log.Fatal("Fatal error while creating local private key output file: ", err)
	}
	defer privateOutFile.Close()

	privateKey := &pem.Block{
		Type:	"PRIVATE KEY",
		Bytes:	x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(privateOutFile, privateKey)
	if err != nil {
		log.Fatal("Fatal error while saving local private key: ", err)
	}

	asn1bytes, err := asn1.Marshal(key.PublicKey)
	if err != nil {
		log.Fatal("Fatal error while calculating local public key: ", err)
	}

	publicOutFile, err := os.Create("keys/local.pub")
	if err != nil {
		log.Fatal("Fatal error while creating local public key output file: ", err)
	}
	defer publicOutFile.Close()

	publicKey := &pem.Block{
		Type:	"PUBLIC KEY",
		Bytes:	asn1bytes,
	}

	err = pem.Encode(publicOutFile, publicKey)
	if err != nil {
		log.Fatal("Fatal error while saving local public key: ", err)
	}

	log.Print("Public key generated:", publicKey)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	if _, err := os.Stat("keys/external.pub"); os.IsNotExist(err) {
		log.Fatal("You must put external public key in ./keys/external.pub")
	}

	if _, err := os.Stat("keys/local"); os.IsNotExist(err) {
		log.Print("Local keypair does not exist, generating now.")
		generateAndSaveKeypair()
	}

	service := microservice.NewPanelMicroservice(":"+port, loadExternal("keys/external.pub"), loadPrivate("keys/local"), "", "")
	service.Endpoints["/test"] = handleTest
	service.Start()
}
