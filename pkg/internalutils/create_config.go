package internalutils

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func CreateConfig() error {
	validatePath := func(filePath string) error {
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) { return nil }
		return err
	}

	prompt := promptui.Prompt{
		Label:    "Path to save the config file",
		Validate: validatePath,
	}

	filePath, err := prompt.Run()

	if err != nil {
		return err
	}

	fileInfo, _ := os.Stat(filePath)
	if fileInfo != nil {
		confirmationPrompt := promptui.Prompt{
			Label: "File exists. Confirm overwrite",
			IsConfirm: true,
		}
		confirm, _ := confirmationPrompt.Run()
		if !strings.EqualFold(confirm, "y") {
			return fmt.Errorf("file exists, chose a different path")
		}
	}

	var queueRegexps []string

	for {
		validateRegex := func(regex string) error {
			_, err := regexp.Compile(regex)
			return err
		}

		regexPrompt := promptui.Prompt{
			Label: "Add a valid regexp to filter queues by name. Leave empty to exit",
			Validate: validateRegex,
		}

		regex, err := regexPrompt.Run()
		if err != nil {
			fmt.Errorf("regexp not valid: %v", err)
			continue
		}

		if strings.TrimSpace(regex) == "" {
			break
		}

		queueRegexps = append(queueRegexps, regex)
	}

	if len(queueRegexps) > 0 {
		fmt.Println("Writing config file: ", filePath)
		viper.Set("filters", map[string][]string{})
		viper.Set("filters.queues", queueRegexps)
		if err := viper.WriteConfigAs(filePath); err != nil {
			return err
		}

		fmt.Println("Config created:")
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	} else {
		fmt.Println("Config file was not created")
	}

	return nil
}
