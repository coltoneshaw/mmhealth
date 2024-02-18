package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	generate "github.com/coltoneshaw/healthcheck/healthcli/generate"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a check to the yaml file",
	Long:  "Interactive dialog to add a check to the yaml file before building the check. ",
	RunE:  addCmdF,
}

func init() {

	RootCmd.AddCommand(
		AddCmd,
	)
}

var qs = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "What is the name of this check?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "group",
		Prompt: &survey.Select{
			Message: "Choose a check group:",
			Options: []string{"config", "packet", "mattermostLog", "notificationLog", "plugins"},
		},
		Validate: survey.Required,
	},
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose a check type:",
			Options: []string{"proactive", "health", "adoption"},
		},
		Validate: survey.Required,
	},
	{
		Name: "severity",
		Prompt: &survey.Select{
			Message: "Choose a check severity:",
			Options: []string{"urgent", "high", "medium", "low"},
			Default: "medium",
		},
		Validate: survey.Required,
	},
	{
		Name:     "description",
		Prompt:   &survey.Input{Message: "What is the description of this check?"},
		Validate: survey.Required,
	},
	{
		Name:     "pass",
		Prompt:   &survey.Input{Message: "What is the pass message?"},
		Validate: survey.Required,
	},
	{
		Name:     "fail",
		Prompt:   &survey.Input{Message: "What is the fail message?"},
		Validate: survey.Required,
	},
	{
		Name: "ignore",
		Prompt: &survey.Input{
			Message: "What is the ignore message? (Optional)",
			Help:    "If you don't want to show anything, just press enter",
		},
	},
}

func addCmdF(cmd *cobra.Command, args []string) error {
	answers := struct {
		Name        string
		Type        string
		Group       string
		Severity    string
		Description string
		Pass        string
		Fail        string
		Ignore      string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return errors.Wrap(err, "Failed to ask questions")
	}

	checks, err := readChecksFile()
	if err != nil {
		return errors.Wrap(err, "Failed to read checks file")
	}

	newKey := generateCheckKey(answers.Type, checks)

	newCheck := generate.Check{
		Name:        answers.Name,
		Result:      generate.Result{Pass: answers.Pass, Fail: answers.Fail, Ignore: answers.Ignore},
		Description: answers.Description,
		Severity:    answers.Severity,
		Type:        generate.CheckType(answers.Type),
	}

	switch answers.Group {
	case "config":
		checks.Config[newKey] = newCheck
		checks.Config = sortGroup(checks.Config)
	case "packet":
		checks.Packet[newKey] = newCheck
		checks.Config = sortGroup(checks.Config)
	case "mattermostLog":
		checks.MattermostLog[newKey] = newCheck
		checks.Config = sortGroup(checks.Config)
	case "notificationLog":
		checks.NotificationLog[newKey] = newCheck
		checks.Config = sortGroup(checks.Config)
	case "plugins":
		checks.Plugins[newKey] = newCheck
		checks.Config = sortGroup(checks.Config)
	}

	// Marshal the Config struct back into YAML
	return storeChecksFile(checks)

}

// parses the existing yaml file and finds the highest existing value and returns the next value
func generateCheckKey(checkType string, checks generate.Checks) string {
	prefix := string(checkType[0])
	highest := 0

	for key := range checks.Config {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.MattermostLog {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.NotificationLog {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.Plugins {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	for key := range checks.Packet {
		if strings.HasPrefix(key, prefix) {
			num, err := strconv.Atoi(key[1:])
			if err == nil && num > highest {
				highest = num
			}
		}
	}
	return fmt.Sprintf("%s%03d", prefix, highest+1)
}

func sortGroup(checks map[string]generate.Check) map[string]generate.Check {
	var keys []string
	for k := range checks {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(keys[i], keys[j]) < 0
	})

	// Create a new sorted map
	sortedChecks := make(map[string]generate.Check)
	for _, k := range keys {
		sortedChecks[k] = checks[k]
	}

	// Replace the 'config' group with the sorted map
	return sortedChecks
}

func readChecksFile() (generate.Checks, error) {
	data, err := os.ReadFile("checks.yaml")
	if err != nil {
		return generate.Checks{}, errors.Wrap(err, "failed to read file")
	}

	var checks generate.Checks
	err = yaml.Unmarshal(data, &checks)
	if err != nil {
		return generate.Checks{}, errors.Wrap(err, "Failed to unmarshal file")
	}

	return checks, nil
}

func storeChecksFile(checks generate.Checks) error {
	data, err := yaml.Marshal(&checks)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal checks file")

	}
	return os.WriteFile("checks.yaml", data, 0644)
}
