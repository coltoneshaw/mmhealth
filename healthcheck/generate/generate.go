package processpacket

import (
	"os"
	"sort"

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

func (p *ProcessPacket) ProcessPacket(packet PacketData) (CheckResults, error) {

	// input file
	checks, err := p.readChecksFile()
	if err != nil {
		return CheckResults{}, err
	}

	p.Checks = checks

	p.Results.Config = p.configChecks(packet.Config)
	p.Results.MattermostLog = p.logChecks(packet.Logs)

	return p.Results, nil
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
