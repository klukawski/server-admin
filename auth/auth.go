package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
)

type AuthConfig struct {
	Validator               *jwt.Validator
	External                *rsa.PublicKey
	Local                   *rsa.PrivateKey
	TlsCertFile, TlsKeyFile string
}

var Config = &AuthConfig{}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := jws.ParseJWTFromRequest(r)
		if ok != nil || token.Validate(Config.External, crypto.SigningMethodRS256, Config.Validator) != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoadPrivate(filename string) *rsa.PrivateKey {
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

func LoadExternal(filename string) *rsa.PublicKey {
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
