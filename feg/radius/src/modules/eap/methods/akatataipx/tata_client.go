package akatataipx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TataClient A client to the TATA IPX
type TataClient struct {
	username string
	password string
	url      string
}

// NewClient ...
func NewClient(username string, password string, url string) (TataClient, error) {
	return TataClient{
		username: username,
		password: password,
		url:      url,
	}, nil
}

// MapsiRequest ...
type MapsiRequest struct {
	Action       string `json:"action_code"`
	UserID       string `json:"userid"`
	Password     string `json:"password"`
	IMSI         string `json:"imsi"`
	TransationID string `json:"transactionid,omitempty"`
}

// MapsiResponse ...
type MapsiResponse struct {
	Action           string      `json:"action_code"`
	ErrorCode        interface{} `json:"errorcode"`
	ErrorString      string      `json:"errorstr"`
	IMSI             string      `json:"imsi"`
	TransationID     string      `json:"transactionid,omitempty"`
	AuthResponseType string      `json:"auth_resp_type,omitempty"`
	AuthTriplets     struct {
		RAND string `json:"randvector"`
		SRES string `json:"sresvector"`
		KC   string `json:"kcvector"`
	} `json:"auth_triplets,omitempty"`
	AuthQuintuplets struct {
		AUTN string `json:"AUTN"`
		CK   string `json:"Ck"`
		IK   string `json:"Ik"`
		RAND string `json:"Rand"`
		XRES string `json:"Xres"`
	} `json:"auth_quintuplets,omitempty"`
}

// Mapsi perform a Mapsi call to TATA IPX interface
func (c TataClient) Mapsi(imsi string, transactionID string) (*MapsiResponse, error) {
	req := MapsiRequest{
		Action:       "mapsai",
		UserID:       c.username,
		Password:     c.password,
		IMSI:         imsi,
		TransationID: transactionID,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resBytes, err := c.send(reqBytes)
	if err != nil {
		return nil, err
	}

	var response MapsiResponse
	err = json.Unmarshal(resBytes, &response)
	if err != nil {
		return nil, err
	}

	// // Hardcoded value for dev/int phases
	// response := MapsiResponse{
	// 	Action:       "mapsai",
	// 	ErrorCode:    0,
	// 	ErrorString:  "SUCCESS",
	// 	IMSI:         "425020699331003",
	// 	TransationID: transactionID,
	// 	AuthTriplets: struct {
	// 		RAND string `json:"randvector"`
	// 		SRES string `json:"sresvector"`
	// 		KC   string `json:"kcvector"`
	// 	}{},
	// 	AuthQuintuplets: struct {
	// 		AUTN string `json:"AUTN"`
	// 		CK   string `json:"Ck"`
	// 		IK   string `json:"Ik"`
	// 		RAND string `json:"Rand"`
	// 		XRES string `json:"Xres"`
	// 	}{
	// 		AUTN: "767D2EE74D0A00009260A1DB0FEE9BA4",
	// 		CK:   "643E97297FA7FA8AA0B2623C432E30B1",
	// 		IK:   "595F559AE486BBECA2D27455E24F4C06",
	// 		RAND: "BA80F6326CBCA976100F677983507B04",
	// 		XRES: "B6DD36C6C30CEE66",
	// 	},
	// }

	return &response, nil
}

// send Sends a POST request to the TATA IPX endpoint and returns the content as string
func (c TataClient) send(body []byte) ([]byte, error) {
	// Create HTTP request
	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Do send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle error codes with clear error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"got status %d from HTTP call with status '%s'",
			resp.StatusCode,
			resp.Status,
		)
	}

	// Read response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Return
	return responseBody, nil
}
