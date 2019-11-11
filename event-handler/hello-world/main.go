package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context/ctxhttp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	// Env is the environment: local or prod
	Env string

	// AppConfigPath is the location where SSM app config is stored
	AppConfigPath string

	// Clients to look for in events.
	Clients []string

	// OpenhabURL is the base URL of the Openhab item
	OpenhabURL string

	// OpenhabAuth is the based64 encoded content for the Authorization header
	OpenhabAuth string

	// Secret is a shared secret with Meraki alerts webhook integration
	Secret string
)

// var sess = session.Must(session.NewSessionWithOptions(session.Options{
// 	SharedConfigState: session.SharedConfigEnable,
// }))
// var ssmClient = ssm.New(sess)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	xray.Configure(xray.Config{LogLevel: "info"})

	fmt.Printf("%s: \n%s\n", request.HTTPMethod, strings.ReplaceAll(request.Body, "\n", ""))
	setupEnvironment(ctx)

	alert, err := parseAlert([]byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Can't parse request: %v", err),
			StatusCode: 500,
		}, nil
	}
	fmt.Printf("Type: %s\n", alert.AlertType)

	if !authorize(alert) {
		return events.APIGatewayProxyResponse{
			Body:       "Unauthorized - bad secret",
			StatusCode: 401,
		}, nil
	}

	switch alert.AlertType {
	case clientConnectivityChanged:
		data := alert.AlertData.(ClientConnectivityChanged)
		connected, _ := strconv.ParseBool(data.Connected)
		// Print log
		state := "ON"
		if !connected {
			state = "OFF"
		}
		fmt.Printf("%s [%s]: %s\n", data.ClientName, data.MAC, state)
		// Post to Openhab
		bodyReader := strings.NewReader(state)
		trimmedMac := strings.ToUpper(data.MAC)
		if !contains(Clients, trimmedMac) {
			fmt.Printf("Not interested in %s\n", trimmedMac)
			break
		}
		trimmedMac = strings.ReplaceAll(trimmedMac, ":", "")
		url := fmt.Sprintf("%s/rest/items/mer_%s", OpenhabURL, trimmedMac)
		fmt.Printf("Posting to %s\n", url)

		req, err := http.NewRequest("POST", url, bodyReader)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Can't create request: %v", err),
				StatusCode: 500,
			}, nil
		}
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", OpenhabAuth))
		// TODO: only do this in non-prod
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		ctx, segment := xray.BeginSubsegment(ctx, "openhab request")
		resp, err := ctxhttp.Do(ctx, xray.Client(nil), req)
		segment.Close(err)

		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Can't send request: %v", err),
				StatusCode: 500,
			}, nil
		}
		if resp.StatusCode != 200 {
			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Error response: %v", resp.Status),
				StatusCode: resp.StatusCode,
			}, nil
		}
		fmt.Printf("Openhab item %s set to %s\n", trimmedMac, state)
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello!\nAlertID: %s\nRequestID: %s\n", alert.AlertID, request.RequestContext.RequestID),
		StatusCode: 200,
	}, nil
}

func authorize(alert Alert) bool {
	if Secret != alert.SharedSecret {
		return false
	}
	return true
}

// help function to check if element is in array
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// load environment and configuration from SSM
func setupEnvironment(ctx context.Context) {
	Env = os.Getenv("ENV")
	if Env == "" {
		fmt.Printf("ENV not configured\n")
	}
	AppConfigPath = os.Getenv("APP_CONFIG_PATH")
	if AppConfigPath == "" {
		fmt.Printf("APP_CONFIG_PATH not configured'n")
	}

	params, err := getParamsFromSSM(ctx)
	if err != nil {
		fmt.Printf("Failed to get from SSM: %v\n", err)
	}

	Clients = params.Clients
	fmt.Printf("Loaded %d devices\n", len(Clients))

	OpenhabURL = params.OpenhabURL
	authString := fmt.Sprintf("%s:%s", params.OpenhabUser, params.OpenhabPwd)
	OpenhabAuth = base64.StdEncoding.EncodeToString([]byte(authString))

	Secret = params.Secret
}

// Parameters is used to parse SSM Parameters
type Parameters struct {
	OpenhabURL  string   `json:"openhab_url"`
	OpenhabUser string   `json:"openhab_user"`
	OpenhabPwd  string   `json:"openhab_pwd"`
	Secret      string   `json:"secret"`
	Clients     []string `json:"clients"`
}

// /merakioh/local
// /merakioh/prod
func getParamsFromSSM(ctx context.Context) (*Parameters, error) {
	keyname := fmt.Sprintf("/%s/%s", AppConfigPath, Env)
	fmt.Printf("get params from %s\n", keyname)

	sess := session.Must(session.NewSession())
	ssmClient := ssm.New(sess)
	xray.AWS(ssmClient.Client)

	withDecryption := false
	if Env == "prod" {
		withDecryption = true
	}
	param, err := ssmClient.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})

	if err != nil {
		return &Parameters{}, err
	}
	if Env == "local" { // don't print sensitive data in prod
		fmt.Printf("Got SSM param %s = %v\n", keyname, *param.Parameter.Value)
	}
	var params *Parameters
	if err := json.Unmarshal([]byte(*param.Parameter.Value), &params); err != nil {
		return &Parameters{}, err
	}
	missingParams := []string{}
	if params.OpenhabURL == "" {
		missingParams = append(missingParams, "OpenhabURL")
	}
	if params.OpenhabUser == "" {
		missingParams = append(missingParams, "OpenhabUser")
	}
	if params.OpenhabPwd == "" {
		missingParams = append(missingParams, "OpenhabPwd")
	}
	if params.Secret == "" {
		missingParams = append(missingParams, "Secret")
	}
	if len(missingParams) != 0 {
		missing := strings.Join(missingParams, ",")
		return &Parameters{}, fmt.Errorf("Missing SSM parameters: %s", missing)
	}
	return params, nil
}

func init() {
}

func main() {
	lambda.Start(handler)
}
