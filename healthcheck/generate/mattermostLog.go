package processpacket

import (
	"bytes"
)

type LogCheckFunc func(checks map[string]Check) CheckResult

func (p *ProcessPacket) logChecks(logs []byte) (results []CheckResult) {

	checks := map[string]LogCheckFunc{
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

func (p *ProcessPacket) h003(checks map[string]Check) CheckResult {
	check, result := initCheckResult("h003", checks, Pass)

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(p.packet.Logs, []byte("context deadline exceeded")) {
		result.Status = Fail
		result.Result = check.Result.Fail
		return result
	}

	return result
}

func (p *ProcessPacket) h004(checks map[string]Check) CheckResult {
	check, result := initCheckResult("h004", checks, Pass)

	// Check if logs contain "i/o timeout"
	if bytes.Contains(p.packet.Logs, []byte("i/o timeout")) {
		result.Status = Fail
		result.Result = check.Result.Fail
		return result
	}
	return result
}

func (p *ProcessPacket) h005(checks map[string]Check) CheckResult {
	check, result := initCheckResult("h005", checks, Pass)

	if bytes.Contains(p.packet.Logs, []byte("Error while creating session for user access token")) {
		// If it does, return a CheckResult with the specified values
		result.Status = Fail
		result.Result = check.Result.Fail
		return result
	}

	// If it doesn't, return a default CheckResult
	return result
}
