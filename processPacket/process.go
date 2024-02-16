package processpacket

import (
	"log"
	"os"

	md "github.com/go-spectest/markdown"
	"gopkg.in/yaml.v3"
)

type CheckType string

const (
	Proactive CheckType = "proactive"
	Health    CheckType = "health"
	Adoption  CheckType = "adoption"
)

type CheckStatus string

const (
	Fail   CheckStatus = "🔴"
	Pass   CheckStatus = "🟢"
	Warn   CheckStatus = "🟡"
	Ignore CheckStatus = "-"
)

type CheckSeverity string

const (
	Urgent CheckSeverity = "urgent"
	High   CheckSeverity = "high"
	Medium CheckSeverity = "medium"
	Low    CheckSeverity = "low"
)

type CheckResult struct {
	ID          string
	Name        string
	Result      string
	Type        CheckType
	Description string
	Status      CheckStatus
	Severity    CheckSeverity
}

type Result struct {
	Pass   string `yaml:"pass"`
	Fail   string `yaml:"fail"`
	Ignore string `yaml:"ignore"`
}

type Check struct {
	Name        string    `yaml:"name"`
	Result      Result    `yaml:"result"`
	Description string    `yaml:"description"`
	Severity    string    `yaml:"severity"`
	Type        CheckType `yaml:"type"`
}

type Checks struct {
	Config          map[string]Check `yaml:"config"`
	Packet          map[string]Check `yaml:"packet"`
	MattermostLog   map[string]Check `yaml:"mattermostLog"`
	NotificationLog map[string]Check `yaml:"notificationLog"`
	Plugins         map[string]Check `yaml:"plugins"`
}

type CheckResults struct {
	Config          []CheckResult
	Packet          []CheckResult
	MattermostLog   []CheckResult
	NotificationLog []CheckResult
	Plugins         []CheckResult
}

type ProcessPacket struct {
	Markdown *md.Markdown
	Checks   Checks
	Results  CheckResults
}

func (p *ProcessPacket) ProcessPacket(packet PacketData) error {

	// input file
	checks, err := p.readChecksFile()
	if err != nil {
		return err
	}

	p.Checks = checks

	p.configChecks(packet.Config)
	p.logChecks(packet.Logs)

	err = p.generateReport()

	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessPacket) readChecksFile() (Checks, error) {
	var checks Checks
	data, err := os.ReadFile("checks.yaml")
	if err != nil {
		return checks, err
	}

	err = yaml.Unmarshal(data, &checks)
	if err != nil {
		return checks, err
	}

	return checks, nil
}

func (p *ProcessPacket) generateReport() error {
	file, err := os.Create("report.md")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()
	p.Markdown = md.NewMarkdown(file).
		H1("Mattermost Health Check Report").
		PlainText("This is an auto generated report from the Mattermost Health Check tool.")

	p.buildResultsTable(p.Results.Config, "Configuration Checks")

	p.buildResultsTable(p.Results.MattermostLog, "Mattermost.Log Checks")
	err = p.Markdown.Build()
	if err != nil {
		return err
	}

	return nil
}

func (p *ProcessPacket) buildResultsTable(testResults []CheckResult, title string) {
	resultsToArray := [][]string{}

	for _, result := range testResults {
		resultsToArray = append(
			resultsToArray,
			[]string{
				result.ID,
				string(result.Type),
				string(result.Severity),
				string(result.Status),
				result.Name,
				result.Result,
				result.Description,
			},
		)
	}

	p.Markdown.
		H2(title).
		CustomTable(md.TableSet{
			Header: []string{"ID", "Result", "Severity", "Status", "Name", "Result", "Description"},
			Rows:   resultsToArray,
		}, md.TableOptions{
			AutoWrapText: false,
		})
}
