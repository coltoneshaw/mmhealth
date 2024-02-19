package processpacket

import (
	"bytes"
)

type LogCheckFunc func(logs []byte, checks map[string]Check) CheckResult

func (p *ProcessPacket) logChecks(logs []byte) (results []CheckResult) {

	checks := map[string]LogCheckFunc{
		"h003": p.h003,
		"h004": p.h004,
		"h005": p.h005,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(logs, p.Checks.MattermostLog)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

func (p *ProcessPacket) h003(logs []byte, checks map[string]Check) CheckResult {
	check, result := initCheckResult("h003", checks, Pass)

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("context deadline exceeded")) {
		result.Status = Fail
		result.Result = check.Result.Fail
	}

	return result
}

func (p *ProcessPacket) h004(logs []byte, checks map[string]Check) CheckResult {
	check, result := initCheckResult("h004", checks, Pass)

	// Check if logs contain "i/o timeout"
	if bytes.Contains(logs, []byte("i/o timeout")) {
		result.Status = Fail
		result.Result = check.Result.Fail
	}
	return result
}

func (p *ProcessPacket) h005(logs []byte, checks map[string]Check) CheckResult {
	check, result := initCheckResult("h005", checks, Pass)

	if bytes.Contains(logs, []byte("Error while creating session for user access token")) {
		// If it does, return a CheckResult with the specified values
		result.Status = Fail
		result.Result = check.Result.Fail
	}

	// If it doesn't, return a default CheckResult
	return result
}
