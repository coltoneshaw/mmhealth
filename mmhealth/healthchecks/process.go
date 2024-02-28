package healthchecks

import (
	"sort"

	"github.com/coltoneshaw/mmhealth/mmhealth/files"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

type CheckResult struct {
	ID          string
	Name        string
	Result      string
	Type        types.CheckType
	Description string
	Status      types.CheckStatus
	Severity    types.CheckSeverity
}

type CheckResults struct {
	Config          []CheckResult
	Packet          []CheckResult
	MattermostLog   []CheckResult
	NotificationLog []CheckResult
	Plugins         []CheckResult
	Environment     []CheckResult
	TopLogs         []types.TopLogs
}

type ProcessPacket struct {
	Checks  types.ChecksFile
	Results CheckResults
	Config  types.ConfigFile
	packet  types.PacketData
}

const (
	Pass   = types.Pass
	Fail   = types.Fail
	Warn   = types.Warn
	Ignore = types.Ignore
)

type CheckFunc func(checks map[string]types.Check) CheckResult

func (p *ProcessPacket) ProcessPacket(packet types.PacketData) (CheckResults, error) {

	p.packet = packet
	// input file
	checksFile, err := files.ReadChecksFile()
	if err != nil {
		return CheckResults{}, err
	}

	p.Checks = checksFile

	configFile, err := files.ReadConfigFile()
	if err != nil {
		return CheckResults{}, err
	}

	p.Config = configFile

	p.Results.Config = p.configChecks(packet.Config)
	p.Results.MattermostLog = p.logChecks()
	p.Results.Environment = p.environmentChecks()

	p.Results.TopLogs = p.topLogs()

	return p.Results, nil
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

func initCheckResult(id string, checks map[string]types.Check, defaultState types.CheckStatus) (types.Check, CheckResult) {
	check := checks[id]

	results := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Description: check.Description,
		Severity:    types.CheckSeverity(check.Severity),
	}

	switch defaultState {
	case types.Fail:
		results.Result = check.Result.Fail
		results.Status = types.Fail

		// Adoption / Proactive checks are not considered fails
		if check.Type == types.Adoption || check.Type == types.Proactive {
			results.Status = types.Warn
		}

	case types.Warn:
		results.Result = check.Result.Fail
		results.Status = types.Warn
	case types.Ignore:
		results.Result = ""
		results.Status = types.Ignore
	case types.Pass:
		results.Result = check.Result.Pass
		results.Status = types.Pass
	}

	return check, results
}
