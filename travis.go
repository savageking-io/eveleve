package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Travis struct {
	conf *TravisConfig
}

type ConfigKey struct {
	Config struct {
		Host        string `json:"host"`
		ShortenHost string `json:"shorten_host"`
		Assets      struct {
			Host string `json:"host"`
		} `json:"assets"`
		Pusher struct {
			Key string `json:"key"`
		} `json:"pusher"`
		Github struct {
			APIURL string   `json:"api_url"`
			Scopes []string `json:"scopes"`
		} `json:"github"`
		Notifications struct {
			Webhook struct {
				PublicKey string `json:"public_key"`
			} `json:"webhook"`
		} `json:"notifications"`
	} `json:"config"`
}

func (t *Travis) Init(config *TravisConfig) error {
	log.Infof("Initializing Travis CI")
	if config == nil {
		return fmt.Errorf("nil travis config")
	}
	t.conf = config

	return nil
}

func (t *Travis) Run() error {
	log.Infof("Starting Travis Listener")
	http.HandleFunc(t.conf.URI, t.Handle)
	return http.ListenAndServe(fmt.Sprintf(":%d", t.conf.Port), nil)
}

func (t *Travis) Handle(w http.ResponseWriter, r *http.Request) {
	log.Infof("New webhook call from Travis")
	log.Debugf("%+v", r)
	key, err := t.TravisPublicKey()
	if err != nil {
		t.RespondWithError(w, err.Error())
		return
	}
	signature, err := t.PayloadSignature(r)
	if err != nil {
		t.RespondWithError(w, err.Error())
		return
	}
	payload := t.PayloadDigest(r.FormValue("payload"))

	err = rsa.VerifyPKCS1v15(key, crypto.SHA1, payload, signature)

	if err != nil {
		t.RespondWithError(w, fmt.Errorf("unauthorized payload").Error())
		return
	}

	log.Infof("Payload: %s", payload)
	t.RespondWithSuccess(w, "payload verified")
}

func (t *Travis) TravisPublicKey() (*rsa.PublicKey, error) {
	log.Debug("Requesting Travis's Public Key")
	response, err := http.Get(t.conf.API)

	if err != nil {
		log.Errorf("Couldn't retrieve public key: %s", err.Error())
		return nil, fmt.Errorf("cannot fetch travis public key")
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var c ConfigKey
	err = decoder.Decode(&c)
	if err != nil {
		log.Errorf("Couldn't decode public key: %s", err.Error())
		return nil, fmt.Errorf("cannot decode travis public key")
	}

	key, err := t.parsePublicKey(c.Config.Notifications.Webhook.PublicKey)
	if err != nil {
		log.Errorf("Failed to parse public key: %s", err.Error())
		return nil, err
	}

	return key, nil
}

func (t *Travis) parsePublicKey(key string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(key))

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid public key")
	}

	return publicKey.(*rsa.PublicKey), nil
}

func (t *Travis) RespondWithError(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	message := fmt.Sprintf("{\"message\": \"%s\"}", m)
	w.Write([]byte(message))
}

func (t *Travis) RespondWithSuccess(w http.ResponseWriter, m string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	message := fmt.Sprintf("{\"message\": \"%s\"}", m)
	w.Write([]byte(message))
}

func (t *Travis) PayloadSignature(r *http.Request) ([]byte, error) {
	signature := r.Header.Get("Signature")
	b64, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Errorf("Couldn't decode signature: %s", err.Error())
		return nil, fmt.Errorf("cannot decode signature")
	}
	return b64, nil
}

func (t *Travis) PayloadDigest(payload string) []byte {
	hash := sha1.New()
	hash.Write([]byte(payload))
	return hash.Sum(nil)
}
