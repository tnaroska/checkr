package config

import (
	"errors"
	"fmt"
	"github.com/analogj/checkr/pkg/utils"
	"github.com/spf13/viper"
)

// When initializing this class the following methods must be called:
// Config.New
// Config.Init
// This is done automatically when created via the Factory.
type configuration struct {
	*viper.Viper
}

//Viper uses the following precedence order. Each item takes precedence over the item below it:
// explicit call to Set
// flag
// env
// config
// key/value store
// default

func (c *configuration) Init() error {
	c.Viper = viper.New()
	//set defaults
	c.SetDefault("details_url", "")

	// Base URL for API requests. Defaults to the public GitHub API, but can be
	// set to a domain endpoint to use with GitHub Enterprise. BaseURL should
	// always be specified with a trailing slash.
	c.SetDefault("base_url", "https://api.github.com/")

	c.SetDefault("org", "")
	c.SetDefault("repo", "")
	c.SetDefault("pr", "")
	c.SetDefault("sha", "")

	c.AutomaticEnv()
	c.SetEnvPrefix("GHCS")

	if c.IsSet("private_key_path") {
		privateKeyPath, err := utils.ExpandPath(c.GetString("private_key_path"))
		if err != nil {
			return err
		}
		if !utils.FileExists(privateKeyPath) {
			return errors.New("GHCS_PRIVATE_KEY_PATH is invalid")
		}
		c.Set("private_key_path", privateKeyPath)
	}

	fmt.Printf("%s", c.AllSettings())
	return nil

}

// This function ensures that the merged config works correctly.
func (c *configuration) ValidateConfig() error {

	if !c.IsSet("pr") && !c.IsSet("sha") {
		return errors.New("Required flag/environmental variable \"pr\" or \"sha\" is not set")
	}

	if !c.IsSet("app_id") {
		return errors.New("GHCS_APP_ID is required")
	}

	if !c.IsSet("private_key") && !c.IsSet("private_key_path") {
		return errors.New("GHCS_PRIVATE_KEY or GHCS_PRIVATE_KEY_PATH is required")
	}

	baseUrl := c.GetString("base_url")
	if !(baseUrl[len(baseUrl)-1:] == "/") {
		return errors.New("GHCS_BASE_URL must include trailing slash '/'")
	}

	return nil
}
