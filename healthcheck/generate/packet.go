package processpacket

// import "github.com/mattermost/mattermost/server/public/model"

// type PacketCheckFunc func(packet model.SupportPacket, checks map[string]Check) CheckResult

// func (p *ProcessPacket) packetChecks(logs []byte) {

// 	checks := map[string]LogCheckFunc{
// 		"h003": h003,
// 		"h004": h004,
// 		"h005": h005,
// 	}
// 	testResults := []CheckResult{}

// 	for id, check := range checks {
// 		result := check(logs, p.Checks.MattermostLog)
// 		result.ID = id
// 		testResults = append(testResults, result)
// 	}

// 	p.Results.MattermostLog = testResults
// }

// func initChecks(id string, checks map[string]Check) CheckResult {
// 	check := checks[id]

// 	results := CheckResult{
// 		Name:        check.Name,
// 		Result:      check.Result.Pass,
// 		Type:        check.Type,
// 		Description: check.Description,
// 		Status:      Pass,
// 		Severity:    CheckSeverity(check.Severity),
// 	}

// 	return results
// }

// // Server Version check
// func h006(logs []byte, checks map[string]Check) CheckResult {

// 	check := initChecks("h006", checks)

// 	// if version >= v9.2 OR is v8.1

// 	// next month if version >= v9.4 OR is v8.1
// 	// stays on v8.1 until may 16th then you change the check to v9.5
// 	// pattern is that one
// 	// need to have some kind of support matrix here

// 	// check if the semvar is within the group
// 	// If it doesn't, return a default CheckResult
// 	return results
// }
