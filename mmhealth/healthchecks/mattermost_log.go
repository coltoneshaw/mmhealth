package healthchecks

import (
	"strings"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func (p *ProcessPacket) logChecks() (results []CheckResult) {

	checks := map[string]CheckFunc{
		"h003": p.h003,
		"h004": p.h004,
		"h005": p.h005,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(p.Checks.MattermostLog)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

func (p *ProcessPacket) h003(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h003", checks, Pass)

	for _, log := range p.packet.Logs {
		if strings.Contains(log.Msg, "context deadline exceeded") {
			result.Status = Fail
			result.Result = check.Result.Fail
			return result
		}
	}

	return result
}

func (p *ProcessPacket) h004(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h004", checks, Pass)

	for _, log := range p.packet.Logs {
		if strings.Contains(log.Msg, "i/o timeout") {
			result.Status = Fail
			result.Result = check.Result.Fail
			return result
		}
	}
	return result
}

func (p *ProcessPacket) h005(checks map[string]types.Check) CheckResult {
	check, result := initCheckResult("h005", checks, Pass)

	for _, log := range p.packet.Logs {
		if strings.Contains(log.Msg, "Error while creating session for user access token") {
			result.Status = Fail
			result.Result = check.Result.Fail
			return result
		}
	}

	return result
}
