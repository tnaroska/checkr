package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/analogj/checkr/pkg/config"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v50/github"
)

func GetJwtClient(appConfig config.Interface) (*github.Client, error) {
	// Wrap the shared transport for use with defined application and installation IDs

	appId := appConfig.GetInt64("app_id")
	var jwtTransport *ghinstallation.AppsTransport
	var err error
	if appConfig.IsSet("private_key_path") {
		jwtTransport, err = ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, appId, appConfig.GetString("private_key_path"))
	} else if appConfig.IsSet("private_key_base64") {

		decodedPrivateKey, err := base64.StdEncoding.DecodeString(appConfig.GetString("private_key_base64"))
		if err != nil {
			fmt.Printf("err: %s", err)
			return nil, err
		}

		jwtTransport, err = ghinstallation.NewAppsTransport(http.DefaultTransport, appId, []byte(decodedPrivateKey))
	}
	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	jwtTransport.BaseURL = strings.TrimSuffix(appConfig.GetString("base_url"), "/")

	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	// Use installation transport with jwtClient
	// NewClient returns a new GitHub API jwtClient.
	// If a nil httpClient is provided, http.DefaultClient will be used. To use API methods which require authentication,
	// provide an http.Client that will perform the authentication for you.
	jwtClient := github.NewClient(&http.Client{Transport: jwtTransport})

	jwtClient.BaseURL, err = url.Parse(appConfig.GetString("base_url"))
	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	return jwtClient, nil
}

func GetAppClient(appConfig config.Interface, installationId int64) (*github.Client, error) {

	appId := appConfig.GetInt64("app_id")
	var appTransport *ghinstallation.Transport
	var err error
	if appConfig.IsSet("private_key_path") {
		appTransport, err = ghinstallation.NewKeyFromFile(http.DefaultTransport, appId, installationId, appConfig.GetString("private_key_path"))
	} else if appConfig.IsSet("private_key_base64") {
		decodedPrivateKey, err := base64.StdEncoding.DecodeString(appConfig.GetString("private_key_base64"))
		if err != nil {
			fmt.Printf("err: %s", err)
			return nil, err
		}

		appTransport, err = ghinstallation.New(http.DefaultTransport, appId, installationId, []byte(decodedPrivateKey))
	}
	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	appTransport.BaseURL = strings.TrimSuffix(appConfig.GetString("base_url"), "/")

	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err

	}
	appClient := github.NewClient(&http.Client{Transport: appTransport})

	appClient.BaseURL, err = url.Parse(appConfig.GetString("base_url"))

	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	//access_token, err := appTransport.Token()
	//if err != nil {
	//	fmt.Printf("err: %s", err)
	//	return nil, err
	//}
	//
	//fmt.Printf("Installation access token: %s", access_token)
	return appClient, nil
}
