package healthchecks

import (
	"slices"
	"sort"

	"github.com/coltoneshaw/mmhealth/mmhealth"
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

type TopLogs struct {
	Count  int
	Caller string
	Msg    string
	Level  string
}

type PluginResults struct {
	PluginID          string
	LatestVersion     string
	LatestReleaseDate string
	InstalledVersion  string
	PluginName        string
	PluginURL         string
	Active            bool
	IsUpdated         bool
	SupportLevel      string
}

type CheckResults struct {
	Config          []CheckResult
	Packet          []CheckResult
	MattermostLog   []CheckResult
	NotificationLog []CheckResult
	Plugins         []PluginResults
	Environment     []CheckResult
	TopLogs         []TopLogs
}

type ProcessPacket struct {
	Checks  types.ChecksFile
	Results CheckResults
	Config  types.ConfigFile
	packet  types.PacketData
	log     func(a ...any)
}

const (
	Pass   = types.StatusPass
	Fail   = types.StatusFail
	Warn   = types.StatusWarn
	Ignore = types.StatusIgnore
	Error  = types.StatusError
)

type CheckFunc func(checks map[string]types.Check) CheckResult

func (p *ProcessPacket) ProcessPacket(packet types.PacketData) (CheckResults, error) {

	p.packet = packet
	p.log = mmhealth.HandleError
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

	p.Results.Config = p.configChecks()
	p.Results.MattermostLog = p.logChecks()
	p.Results.Environment = p.environmentChecks()

	p.Results.TopLogs = p.topLogs()
	p.Results.Packet = p.packetChecks()

	p.Results.Plugins = p.pluginChecks()

	return p.Results, nil
}
func (p *ProcessPacket) sortResults(testResults []CheckResult) []CheckResult {
	var errorResults []CheckResult
	var failResults []CheckResult
	var warnResults []CheckResult
	var passResults []CheckResult
	var ignoreResults []CheckResult

	for _, result := range testResults {
		switch result.Status {
		case types.StatusError:
			errorResults = append(errorResults, result)
		case types.StatusFail:
			failResults = append(failResults, result)
		case types.StatusWarn:
			warnResults = append(warnResults, result)
		case types.StatusPass:
			passResults = append(passResults, result)
		case types.StatusIgnore:
			ignoreResults = append(ignoreResults, result)
		}
	}

	errorResults = p.sortBySev(errorResults)
	failResults = p.sortBySev(failResults)
	warnResults = p.sortBySev(warnResults)
	passResults = p.sortBySev(passResults)
	ignoreResults = p.sortBySev(ignoreResults)

	return slices.Concat(errorResults, failResults, warnResults, passResults, ignoreResults)
}

func (p *ProcessPacket) sortBySev(testResults []CheckResult) []CheckResult {
	severityOrder := map[types.CheckSeverity]int{
		types.Urgent: 0,
		types.High:   1,
		types.Medium: 2,
		types.Low:    3,
	}

	sort.Slice(testResults, func(i, j int) bool {
		return severityOrder[testResults[i].Severity] < severityOrder[testResults[j].Severity]
	})
	return testResults
}

func initCheckResult(id string, checks map[string]types.Check, defaultState types.CheckStatus) (types.Check, CheckResult) {
	check := checks[id]

	results := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Description: check.Description,
		Severity:    check.Severity,
	}

	switch defaultState {
	case types.StatusFail:
		results.Result = check.Result.Fail
		results.Status = types.StatusFail

		// Adoption / Proactive checks are not considered fails
		if check.Type == types.Adoption || check.Type == types.Proactive {
			results.Status = types.StatusWarn
		}

	case types.StatusWarn:
		results.Result = check.Result.Fail
		results.Status = types.StatusWarn
	case types.StatusIgnore:
		results.Result = ""
		results.Status = types.StatusIgnore
	case types.StatusPass:
		results.Result = check.Result.Pass
		results.Status = types.StatusPass
	case types.StatusError:
		results.Result = check.Result.Error
		results.Status = types.StatusError
	}

	return check, results
}
