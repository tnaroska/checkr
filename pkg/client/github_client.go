package client

import (
	"fmt"
	"github.com/analogj/checkr/pkg/config"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"net/http"
)

func GetJwtClient(appConfig config.Interface) (*github.Client, error) {
	// Wrap the shared transport for use with defined application and installation IDs

	appId := appConfig.GetInt("APP_ID")
	var jwtTransport *ghinstallation.AppsTransport
	var err error
	if appConfig.IsSet("PRIVATE_KEY_PATH") {
		jwtTransport, err = ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, appId, appConfig.GetString("PRIVATE_KEY_PATH"))
	} else if appConfig.IsSet("PRIVATE_KEY") {
		jwtTransport, err = ghinstallation.NewAppsTransport(http.DefaultTransport, appId, []byte(appConfig.GetString("PRIVATE_KEY")))
	}

	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	// Use installation transport with jwtClient
	// NewClient returns a new GitHub API jwtClient.
	// If a nil httpClient is provided, http.DefaultClient will be used. To use API methods which require authentication,
	// provide an http.Client that will perform the authentication for you.
	return github.NewClient(&http.Client{Transport: jwtTransport}), nil

}

func GetAppClient(appConfig config.Interface, installationId int) (*github.Client, error) {

	appId := appConfig.GetInt("APP_ID")
	var appTransport *ghinstallation.Transport
	var err error
	if appConfig.IsSet("PRIVATE_KEY_PATH") {
		appTransport, err = ghinstallation.NewKeyFromFile(http.DefaultTransport, appId, installationId, appConfig.GetString("PRIVATE_KEY_PATH"))
	} else if appConfig.IsSet("PRIVATE_KEY") {
		appTransport, err = ghinstallation.New(http.DefaultTransport, appId, installationId, []byte(appConfig.GetString("PRIVATE_KEY")))
	}

	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err

	}
	appClient := github.NewClient(&http.Client{Transport: appTransport})

	access_token, err := appTransport.Token()
	if err != nil {
		fmt.Printf("err: %s", err)
		return nil, err
	}

	fmt.Printf("Installation access token: %s", access_token)
	return appClient, nil
}
