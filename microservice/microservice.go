package microservice

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"crypto/rsa"
)

type Endpoint func(http.ResponseWriter, *http.Request)

type PanelMicroservice struct {
	server                  http.Server
	private					*rsa.PrivateKey
	external				*rsa.PublicKey
	claims                  *jws.Claims
	tlsCertFile, tlsKeyFile string
	Endpoints               map[string]Endpoint
}

func NewPanelMicroservice(address string, external *rsa.PublicKey, privKey *rsa.PrivateKey, tlsCertFile, tlsKeyFile string) *PanelMicroservice {
	panelMicroservice := &PanelMicroservice{
		server: http.Server{
			Addr: address,
		},
		private:	privKey,
		external:	external,
		claims:     &jws.Claims{},
		tlsCertFile:tlsCertFile,
		tlsKeyFile: tlsKeyFile,
		Endpoints:  map[string]Endpoint{},
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
