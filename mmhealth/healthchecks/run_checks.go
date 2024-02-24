package healthchecks

import (
	"os"
	"sort"

	"github.com/coltoneshaw/mmhealth/mmhealth"
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

type ChecksFile struct {
	Config          map[string]Check `yaml:"config"`
	Packet          map[string]Check `yaml:"packet"`
	Environment     map[string]Check `yaml:"environment"`
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
	Environment     []CheckResult
}

type ProcessPacket struct {
	Checks  ChecksFile
	Results CheckResults
	Config  ConfigFile
	packet  mmhealth.PacketData
}

type ConfigFile struct {
	Versions struct {
		Supported []string `yaml:"supported"`
		ESR       string   `yaml:"esr"`
	} `yaml:"versions"`
}

func (p *ProcessPacket) ProcessPacket(packet mmhealth.PacketData) (CheckResults, error) {

	p.packet = packet
	// input file
	checksFile, err := p.readChecksFile()
	if err != nil {
		return CheckResults{}, err
	}

	p.Checks = checksFile

	configFile, err := p.readConfigFile()
	if err != nil {
		return CheckResults{}, err
	}

	p.Config = configFile

	p.Results.Config = p.configChecks(packet.Config)
	p.Results.MattermostLog = p.logChecks(packet.Logs)
	p.Results.Environment = p.environmentChecks()

	return p.Results, nil
}

func (p *ProcessPacket) readChecksFile() (ChecksFile, error) {
	var checks ChecksFile
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

func (p *ProcessPacket) readConfigFile() (ConfigFile, error) {
	var config ConfigFile
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
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

func initCheckResult(id string, checks map[string]Check, defaultState CheckStatus) (Check, CheckResult) {
	check := checks[id]

	results := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Description: check.Description,
		Severity:    CheckSeverity(check.Severity),
	}

	switch defaultState {
	case Fail:
		results.Result = check.Result.Fail
		results.Status = Fail

		// Adoption / Proactive checks are not considered fails
		if check.Type == Adoption || check.Type == Proactive {
			results.Status = Warn
		}

	case Warn:
		results.Result = check.Result.Fail
		results.Status = Warn
	case Ignore:
		results.Result = ""
		results.Status = Ignore
	case Pass:
		results.Result = check.Result.Pass
		results.Status = Pass
	}

	return check, results
}
