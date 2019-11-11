package main

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	apsCameUp                 = "APs came up"
	apsWentDown               = "APs went down"
	clientConnectivityChanged = "Client connectivity changed"
)

// Alert represents a Meraki alert as send by the webhooks
type Alert struct {
	AlertID          string    `json:"alertId"`
	AlertType        string    `json:"alertType"`
	NetworkID        string    `json:"networkId"`
	NetworkName      string    `json:"networkName"`
	NetworkURL       string    `json:"networkUrl"`
	OccurredAt       time.Time `json:"occurredAt"`
	OrganizationID   string    `json:"organizationId"`
	OrganizationName string    `json:"organizationName"`
	OrganizationURL  string    `json:"organizationUrl"`
	SentAt           time.Time `json:"sentAt"`
	SharedSecret     string    `json:"sharedSecret"`
	Version          string    `json:"version"`

	AlertData interface{} `json:"alertData"`
}

// ClientConnectivityChanged represents AlertData for `clientConnectivityChanged` events
type ClientConnectivityChanged struct {
	MAC        string `json:"mac"`
	IP         string `json:"ip"`
	Connected  string `json:"connected"`
	ClientName string `json:"clientName"`
	ClientURL  string `json:"clientUrl"`
}

func parseAlert(rawAlert []byte) (Alert, error) {
	var alertData json.RawMessage
	alert := Alert{
		AlertData: &alertData,
	}

	if err := json.Unmarshal(rawAlert, &alert); err != nil {
		return Alert{}, err
	}

	if alert.AlertType == "" {
		return Alert{}, errors.New("Parsing failed")
	}

	switch alert.AlertType {
	case clientConnectivityChanged:
		var changedData ClientConnectivityChanged
		if err := json.Unmarshal(alertData, &changedData); err != nil {
			return Alert{}, nil
		}
		alert.AlertData = changedData
	}
	return alert, nil
}
