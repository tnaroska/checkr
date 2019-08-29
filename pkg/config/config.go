package config

import (
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
	c.SetDefault("options.config_dir", "~/.ssh/drawbridge")
	c.SetDefault("options.pem_dir", "~/.ssh/drawbridge/pem")
	c.SetDefault("options.active_config_template", "default")

	c.SetDefault("options.active_custom_templates", []string{})
	c.SetDefault("options.ui_group_priority", []string{"environment", "stack_name", "shard", "shard_type"})
	c.SetDefault("options.ui_question_hidden", []string{})

	//if you want to load a non-standard location system config file (~/drawbridge.yml), use ReadConfig
	c.SetConfigType("yaml")
	//c.SetConfigName("drawbridge")
	//c.AddConfigPath("$HOME/")

	//CLI options will be added via the `Set()` function
	return c.ValidateConfig()

}

// This function ensures that the merged config works correctly.
func (c *configuration) ValidateConfig() error {

	////deserialize Questions
	//questionsMap := map[string]Question{}
	//err := c.UnmarshalKey("questions", &questionsMap)
	//
	//if err != nil {
	//	log.Printf("questions could not be deserialized correctly. %v", err)
	//	return err
	//}
	//
	//for _, v := range questionsMap {
	//
	//	typeContent, ok := v.Schema["type"].(string)
	//	if !ok || len(typeContent) == 0 {
	//		return errors.QuestionSyntaxError("`type` is required for questions")
	//	}
	//}
	//
	//

	return nil
}
