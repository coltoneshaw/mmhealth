package processpacket

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
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
	Fail   CheckStatus = "fail"
	Pass   CheckStatus = "pass"
	Warn   CheckStatus = "warn"
	Ignore CheckStatus = "ignore"
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
	Checks  Checks
	Results CheckResults
}

func (p *ProcessPacket) ProcessPacket(packet PacketData) error {

	// input file
	checks, err := p.readChecksFile()
	if err != nil {
		return err
	}

	p.Checks = checks

	p.Results.Config = p.configChecks(packet.Config)
	p.Results.MattermostLog = p.logChecks(packet.Logs)

	err = p.SaveResultsToFile()
	if err != nil {
		return errors.Wrap(err, "failed to save results to file")
	}
	return nil
}

func (p *ProcessPacket) SaveResultsToFile() error {
	// Convert the results to YAML
	data, err := yaml.Marshal(p.Results)
	if err != nil {
		return errors.Wrap(err, "failed to marshal results to yaml")
	}

	file, err := os.Create(filepath.Join("/files", "report.md"))
	if err != nil {
		return errors.Wrap(err, "failed to create report.md")
	}
	defer file.Close()

	markdown := "---\n" + string(data) + "---\n"

	err = os.WriteFile(filepath.Join("/files", "report.md"), []byte(markdown), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write to report.md")
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

func (p *ProcessPacket) sortResults(testResults []CheckResult) []CheckResult {
	statusOrder := map[string]int{
		"fail":   1,
		"warn":   2,
		"pass":   3,
		"ignore": 4,
	}

	sort.Slice(testResults, func(i, j int) bool {
		return statusOrder[string(testResults[i].Status)] < statusOrder[string(testResults[j].Status)]
	})
	return testResults
}
