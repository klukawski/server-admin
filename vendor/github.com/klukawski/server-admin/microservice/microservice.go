package microservice

import (
	"fmt"
	"net/http"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
)

type Endpoint func(http.ResponseWriter, *http.Request)

type PanelMicroservice struct {
	server                  http.Server
	private                 *rsa.PrivateKey
	external                *rsa.PublicKey
	claims                  *jws.Claims
	tlsCertFile, tlsKeyFile string
	Endpoints               map[string]Endpoint
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
	if err != nil {
		log.Fatal("Fatal error while generating a keypair: ", err)
	}

	privateOutFile, err := os.Create("keys/local")
	if err != nil {
		log.Fatal("Fatal error while creating local private key output file: ", err)
	}
	defer privateOutFile.Close()

	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
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
		Type:  "PUBLIC KEY",
		Bytes: asn1bytes,
	}

	err = pem.Encode(publicOutFile, publicKey)
	if err != nil {
		log.Fatal("Fatal error while saving local public key: ", err)
	}

	log.Print("Public key generated:", publicKey)
}

func NewPanelMicroservice(address, external, privKey, tlsCertFile, tlsKeyFile string) *PanelMicroservice {
	if _, err := os.Stat("keys/external.pub"); os.IsNotExist(err) {
		log.Fatal("You must put external public key in ./keys/external.pub")
	}

	if _, err := os.Stat("keys/local"); os.IsNotExist(err) {
		log.Print("Local keypair does not exist, generating now.")
		generateAndSaveKeypair()
	}
	panelMicroservice := &PanelMicroservice{
		server: http.Server{
			Addr: address,
		},
		private:     loadPrivate(privKey),
		external:    loadExternal(external),
		claims:      &jws.Claims{},
		tlsCertFile: tlsCertFile,
		tlsKeyFile:  tlsKeyFile,
		Endpoints:   map[string]Endpoint{},
	}
	panelMicroservice.claims.SetIssuer("panel")
	panelMicroservice.server.Handler = panelMicroservice
	return panelMicroservice
}

func (panel *PanelMicroservice) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	validator := jws.NewValidator(*panel.claims, time.Minute, time.Minute, nil)
	err = t.Validate(panel.external, crypto.SigningMethodRS256, validator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	endpoint, ok := panel.Endpoints[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Endpoint not found")
		return
	}
	endpoint(w, r)
}

func (panel *PanelMicroservice) Start() {
	panel.server.ListenAndServe() //TLS(tlsCertFile, tlsKeyFile)
}
