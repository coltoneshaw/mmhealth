package processpacket

import (
	"os"

	md "github.com/go-spectest/markdown"
)

type CheckType string

const (
	Proactive CheckType = "proactive"
	Health    CheckType = "health"
	Adoption  CheckType = "adoption"
)

type CheckResult struct {
	Name        string
	Result      string
	Type        CheckType
	Description string
	Status      string
}

func ProcessPacket(packet PacketData) error {

	results := md.NewMarkdown(os.Stdout).
		H1("Mattermost Health Check").
		PlainText("This is an auto generated report from the Mattermost Health Check tool.")

	configChecks(packet.Config, results)
	logChecks(packet.Logs, results)

	err := results.Build()

	if err != nil {
		return err
	}
	return nil
}
