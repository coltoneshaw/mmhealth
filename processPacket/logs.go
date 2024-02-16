package processpacket

import (
	"bytes"
	"fmt"

	md "github.com/go-spectest/markdown"
)

type LogCheckFunc func(logs []byte) CheckResult

func logChecks(logs []byte, results *md.Markdown) {

	checks := []LogCheckFunc{contextDeadlineExceeded, ioTimeout, userTokenError}
	testResults := []CheckResult{}

	for _, check := range checks {
		result := check(logs)
		testResults = append(testResults, result)
	}

	resultsToArray := [][]string{}

	for _, result := range testResults {
		resultsToArray = append(resultsToArray, []string{result.Name, string(result.Type), result.Status, result.Result, result.Description})
	}

	fmt.Println(resultsToArray)
	results.
		H2("Mattermost.log Checks").
		CustomTable(md.TableSet{
			Header: []string{"Name", "Type", "Status", "Result", "Description"},
			Rows:   resultsToArray,
		}, md.TableOptions{
			AutoWrapText: false,
		})
}

func contextDeadlineExceeded(logs []byte) CheckResult {

	results := CheckResult{
		Name:        "context deadline exceeded",
		Result:      "context deadline exceeded not found",
		Type:        Health,
		Description: "The context deadline exceeded error is a common error in Mattermost. It is usually caused by a slow database or a slow network connection. [documentation](https://docs.mattermost.com/install/troubleshooting.html#context-deadline-exceeded)",
		Status:      "🟢",
	}

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("context deadline exceeded")) {
		// If it does, return a CheckResult with the specified values
		results.Status = "🔴"
		results.Result = "context deadline exceeded found"
	}

	// If it doesn't, return a default CheckResult
	return results
}

func ioTimeout(logs []byte) CheckResult {

	results := CheckResult{
		Name:        "i/o timeout",
		Result:      "i/o timeout not found",
		Type:        Health,
		Description: "Further investigation is needed.  Contact your Technical Account Manager for assistance. A common cause of this error is due to connectivity issues.  The root cause can originate from various factors. Depending on the origin of the error, we recommend verifying accessibility of the resource. In some cases, ingress/egress rules might be causing problems, or issues may arise from the nginx configuration.",
		Status:      "🟢",
	}

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("i/o timeout")) {
		// If it does, return a CheckResult with the specified values
		results.Status = "🔴"
		results.Result = "i/o timeout found"
	}

	// If it doesn't, return a default CheckResult
	return results
}

func userTokenError(logs []byte) CheckResult {

	results := CheckResult{
		Name:        "Error while creating session",
		Result:      "Error while creating session for user access token not found",
		Type:        Health,
		Description: "",
		Status:      "🟢",
	}

	// Check if logs contain "context deadline exceeded"
	if bytes.Contains(logs, []byte("Error while creating session for user access token")) {
		// If it does, return a CheckResult with the specified values
		results.Status = "🔴"
		results.Result = "Error while creating session for user access token found"
	}

	// If it doesn't, return a default CheckResult
	return results
}
