package microservice

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
)

type Endpoint func(http.ResponseWriter, *http.Request)

type PanelMicroservice struct {
	server                  http.Server
	key                     []byte
	claims                  *jws.Claims
	tlsCertFile, tlsKeyFile string
	Endpoints               map[string]Endpoint
}

func token() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (panel *PanelMicroservice) handleToken(w http.ResponseWriter, r *http.Request) {
	token, _ := panel.claims.JWTID()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"token\":\"%s\"}", token)
}

func NewPanelMicroservice(address, key, tlsCertFile, tlsKeyFile string) *PanelMicroservice {
	keyBytes, _ := base64.StdEncoding.DecodeString(key)
	panelMicroservice := &PanelMicroservice{
		server: http.Server{
			Addr: address,
		},
		key:         keyBytes,
		claims:      &jws.Claims{},
		tlsCertFile: tlsCertFile,
		tlsKeyFile:  tlsKeyFile,
		Endpoints:   map[string]Endpoint{},
	}
	panelMicroservice.claims.SetIssuer("panel")
	panelMicroservice.claims.SetJWTID(token())
	panelMicroservice.server.Handler = panelMicroservice
	return panelMicroservice
}

func (panel *PanelMicroservice) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/token" {
		panel.handleToken(w, r)
		return
	}

	t, err := jws.ParseJWTFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	validator := jws.NewValidator(*panel.claims, time.Minute, time.Minute, nil)
	err = t.Validate(panel.key, crypto.SigningMethodHS512, validator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	panel.claims.SetJWTID(token())

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
