package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestAlerts(t *testing.T) {
	rawAlert := `
	{
		"version": "0.1",
		"sharedSecret": "verysecret",
		"sentAt": "2019-11-09T08:52:39.258809Z",
		"organizationId": "456954",
		"organizationName": "Tom Deckers",
		"organizationUrl": "https://n207.meraki.com/o/aa9sQd/manage/organization/overview",
		"networkId": "L_665406844943993490",
		"networkName": "Home",
		"networkUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/wired_status",
		"deviceSerial": "Q2MN-5J4X-FLQD",
		"deviceMac": "00:18:0a:3c:83:b0",
		"deviceName": "Gateway",
		"deviceUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/new_wired_status",
		"deviceTags": [],
		"deviceModel": "MX64W",
		"alertId": "679480593797143137",
		"alertType": "Client connectivity changed",
		"occurredAt": "2019-11-09T08:52:17.913000Z",
		"alertData": {
			"mac": "3C:28:6D:29:A7:66",
			"ip": "192.168.2.246",
			"connected": "false",
			"clientName": "Pixel 3 (Tom)",
			"clientUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/usage/list#c=k567cb7"
		}
	}`

	alert, err := parseAlert([]byte(rawAlert))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, alert.Version, "0.1")
	assert.Equal(t, alert.SharedSecret, "verysecret")
	assert.Equal(t, alert.SentAt.Day(), 9)
	assert.Equal(t, alert.OrganizationID, "456954")
	assert.Equal(t, alert.OrganizationName, "Tom Deckers")
	assert.Equal(t, alert.OrganizationURL, "https://n207.meraki.com/o/aa9sQd/manage/organization/overview")
	assert.Equal(t, alert.NetworkID, "L_665406844943993490")
	assert.Equal(t, alert.NetworkName, "Home")
	assert.Equal(t, alert.NetworkURL, "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/wired_status")
	assert.Equal(t, alert.AlertType, clientConnectivityChanged)
	assert.Equal(t, alert.OccurredAt.Day(), 9)
	assert.Equal(t, alert.AlertID, "679480593797143137")

	changedData, ok := alert.AlertData.(ClientConnectivityChanged)
	if !ok {
		t.Fatalf("Failed to convert type\n")
	}
	assert.Equal(t, changedData.ClientName, "Pixel 3 (Tom)")

}
