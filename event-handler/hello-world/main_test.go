package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var (
	req1 = events.APIGatewayProxyRequest{Body: `{    "version": "0.1",    "sharedSecret": "verysecret",    "sentAt": "2019-11-09T08:56:18.503438Z",    "organizationId": "456954",    "organizationName": "Tom Deckers",    "organizationUrl": "https://n207.meraki.com/o/aa9sQd/manage/organization/overview",    "networkId": "L_665406844943993490",    "networkName": "Home",    "networkUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/wired_status",    "deviceSerial": "Q2MN-5J4X-FLQD",    "deviceMac": "00:18:0a:3c:83:b0",    "deviceName": "Gateway",    "deviceUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/nodes/new_wired_status",    "deviceTags": [],    "deviceModel": "MX64W",    "alertId": "679480593797143141",    "alertType": "Client connectivity changed",    "occurredAt": "2019-11-09T08:55:45.422000Z",    "alertData": {        "mac": "3C:28:6D:29:A7:66",        "ip": "192.168.2.246",        "connected": "true",        "clientName": "Pixel 3 (Tom)",        "clientUrl": "https://n207.meraki.com/Home-appliance/n/OOw3Oapd/manage/usage/list#c=k567cb7"    }}`}
)

func setup() {
	os.Setenv("DEVICES", "1")
	os.Setenv("DEVICE1", "3C:28:6D:29:A7:66")
	os.Setenv("SECRET", "verysecret")
	setupEnvironment()
}

func TestHandler(t *testing.T) {
	setup()

	t.Run("Unable to get IP", func(t *testing.T) {

		_, err := handler(events.APIGatewayProxyRequest{})
		if err == nil {
			t.Fatal("Error failed to trigger with an invalid request")
		}
	})

	t.Run("Non 200 Response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))
		defer ts.Close()

		OpenhabURL = ts.URL

		_, err := handler(req1)
		if err == nil {
			t.Fatalf("Error failed to trigger with an invalid HTTP response: %v", err)
		}
	})

	t.Run("Invalid secret", func(t *testing.T) {
		origSecret := Secret
		Secret = Secret + "----"

		_, err := handler(req1)

		Secret = origSecret // Reset for other tests.
		if err == nil {
			t.Fatalf("Error failed to trigger with invalid secret: %v", err)
		}
	})

	t.Run("Verify POST", func(t *testing.T) {
		trimmedMac := strings.ReplaceAll(os.Getenv("DEVICE1"), ":", "")
		urlString := fmt.Sprintf("%s/rest/items/mer_%s", OpenhabURL, trimmedMac)
		url, _ := url.Parse(urlString)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != url.Path {
				t.Fatalf("URL doesn't match: %s != %s", r.URL.Path, url.Path)
			}
			w.WriteHeader(200)
		}))
		defer ts.Close()

		OpenhabURL = ts.URL

		_, err := handler(req1)
		if err != nil {
			t.Fatalf("POST did't verify correctly: %v", err)
		}
	})

	// TODO: can mock ssm, see: https://docs.aws.amazon.com/sdk-for-go/api/service/ssm/

}
