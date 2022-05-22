package lyra

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"imanukula/lyra-client/pkg/lyra/response"
)

const (
	DefaultContext   = "TEST"
	DefaultUser      = "69876357"
	DefaultPassword  = "testpassword_DEMOPRIVATEKEY23G4475zXZQ2UA5x7M"
	DefaultEndpoint  = "https://api.payzen.eu"
	DefaultPublicKey = "69876357:testpublickey_DEMOPUBLICKEY95me92597fd28tGD4r5"
	DefaultHashKey   = "38453613e7f44dc58732bad3dca2bca3"
)

type Client struct {
	context    *context.Context
	httpClient *http.Client

	config EpayncConfig
}

type EpayncConfig struct {
	Mode      string
	User      string
	Password  string
	Endpoint  string
	PublicKey string
	HashKey   string
}

func NewClient() (client *Client) {
	return &Client{
		httpClient: &http.Client{},
		config: EpayncConfig{
			Mode:      DefaultContext,
			User:      DefaultUser,
			Password:  DefaultPassword,
			Endpoint:  DefaultEndpoint,
			PublicKey: DefaultPublicKey,
			HashKey:   DefaultHashKey,
		},
	}
}

func (c *Client) GetEndpoint() string {
	return c.config.Endpoint
}

func (c *Client) GetPublicKey() string {
	return c.config.PublicKey
}

func (c *Client) GetHashKey() string {
	return c.config.HashKey
}

func (c *Client) SetContext(ctx context.Context) {
	c.context = &ctx
}

// Authorization - Génère le token d'autorisation pour les requêtes API vers le serveur Epaync
// docs : https://epaync.nc/doc/fr-FR/rest/V4.0/api/kb/authentication.html
func (c *Client) Authorization() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.config.User, c.config.Password)))
}

// CreatePayment - Créer une transaction
// docs : https://epaync.nc/doc/fr-FR/rest/V4.0/api/playground/Charge/CreatePayment/
func (c *Client) CreatePayment(params interface{}) (res *response.EpayncResponse, err error) {
	requestBody, err := json.Marshal(&params)
	if err != nil {
		log.Panicf("error json marshal err=%v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(*c.context, http.MethodPost, c.config.Endpoint+"/api-payment/V4/Charge/CreatePayment", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Panicf("error, build query http, err=%v", err)
		return
	}

	return c.do(req)
}

// do
func (client *Client) do(req *http.Request) (res *response.EpayncResponse, err error) {

	if client.config.User == "" {
		return nil, errors.New("username is not defined in the SDK")
	}

	if client.config.Password == "" {
		return nil, errors.New("password is not defined in the SDK")
	}

	if client.config.Endpoint == "" {
		return nil, errors.New("REST API endpoint not defined in the SDK")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", client.Authorization()))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return
}
