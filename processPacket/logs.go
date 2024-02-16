package processpacket

import (
	"bytes"
)

type LogCheckFunc func(logs []byte, checks map[string]Check) CheckResult

func (p *ProcessPacket) logChecks(logs []byte) {

	checks := map[string]LogCheckFunc{
		"h003": h003,
		"h004": h004,
		"h005": h005,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(logs, p.Checks.MattermostLog)
		result.ID = id
		testResults = append(testResults, result)
	}

	p.Results.MattermostLog = testResults

}

func h003(logs []byte, checks map[string]Check) CheckResult {
	check := checks["h003"]

	results := CheckResult{
		Name:        check.Name,
		Result:      check.Result.Pass,
		Type:        check.Type,
		Description: check.Description,
		Status:      Pass,
		Severity:    CheckSeverity(check.Severity),
	}

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("context deadline exceeded")) {
		// If it does, return a CheckResult with the specified values
		results.Status = Fail
		results.Result = check.Result.Fail
	}

	// If it doesn't, return a default CheckResult
	return results
}

func h004(logs []byte, checks map[string]Check) CheckResult {
	check := checks["h004"]

	results := CheckResult{
		Name:        check.Name,
		Result:      check.Result.Pass,
		Type:        check.Type,
		Description: check.Description,
		Status:      Pass,
		Severity:    CheckSeverity(check.Severity),
	}
	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("i/o timeout")) {
		// If it does, return a CheckResult with the specified values
		results.Status = Fail
		results.Result = check.Result.Fail
	}

	// If it doesn't, return a default CheckResult
	return results
}

func h005(logs []byte, checks map[string]Check) CheckResult {
	check := checks["h005"]

	results := CheckResult{
		Name:        check.Name,
		Result:      check.Result.Pass,
		Type:        check.Type,
		Description: check.Description,
		Status:      Pass,
		Severity:    CheckSeverity(check.Severity),
	}

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("Error while creating session for user access token")) {
		// If it does, return a CheckResult with the specified values
		results.Status = Fail
		results.Result = check.Result.Fail
	}

	// If it doesn't, return a default CheckResult
	return results
}
