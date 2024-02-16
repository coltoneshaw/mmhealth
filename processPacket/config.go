package processpacket

import (
	"fmt"
	"regexp"

	"github.com/mattermost/mattermost/server/public/model"
)

type ConfigCheckFunc func(config model.Config, checks map[string]Check) CheckResult

func (p *ProcessPacket) configChecks(config model.Config) {

	checks := map[string]ConfigCheckFunc{
		"h001": h001,
		"h002": h002,
		"p001": p001,
		"p002": p002,
		"a002": a002,
		"a001": a001,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(config, p.Checks.Config)
		result.ID = id
		testResults = append(testResults, result)
	}

	p.Results.Config = testResults

	// resultsToArray := [][]string{}

	// for _, result := range testResults {
	// 	resultsToArray = append(resultsToArray, []string{result.ID, string(result.Type), result.Name, string(result.Status), result.Result, result.Description})
	// }

	// fmt.Println(resultsToArray)
	// p.Markdown.
	// 	H2("Configuration Checks").
	// 	CustomTable(md.TableSet{
	// 		Header: []string{"ID", "Type", "Name", "Status", "Result", "Description"},
	// 		Rows:   resultsToArray,
	// 	}, md.TableOptions{
	// 		AutoWrapText: false,
	// 	})
}

//
//
// SERVICE SETTINGS
//
//

func h001(config model.Config, checks map[string]Check) CheckResult {
	check := checks["h001"]
	result := CheckResult{
		Name:        check.Name,
		Result:      check.Result.Fail,
		Type:        check.Type,
		Description: check.Description,
		Status:      Fail,
		Severity:    CheckSeverity(check.Severity),
	}
	if *config.ServiceSettings.SiteURL != "" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}
	return result
}

func a001(config model.Config, checks map[string]Check) CheckResult {
	check := checks["a001"]

	result := CheckResult{
		Name:        check.Name,
		Result:      check.Result.Fail,
		Type:        check.Type,
		Description: check.Description,
		Status:      Warn,
		Severity:    CheckSeverity(check.Severity),
	}

	if *config.ServiceSettings.EnableLinkPreviews {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

func a002(config model.Config, checks map[string]Check) CheckResult {
	check := checks["a002"]

	result := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Result:      check.Result.Fail,
		Description: check.Description,
		Status:      Warn,
		Severity:    CheckSeverity(check.Severity),
	}

	if *config.ServiceSettings.ExtendSessionLengthWithActivity {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

//
//
// SQL SETTINGS
//
//

func p001(config model.Config, checks map[string]Check) CheckResult {
	check := checks["p001"]

	result := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Description: check.Description,
		Status:      Warn,
		Result:      check.Result.Fail,
		Severity:    CheckSeverity(check.Severity),
	}

	ipRegexp := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

	containsIPs := false

	if ipRegexp.MatchString(*config.SqlSettings.DataSource) {
		containsIPs = true
	}

	if len(config.SqlSettings.DataSourceReplicas) > 0 {
		for _, replica := range config.SqlSettings.DataSourceReplicas {
			if ipRegexp.MatchString(replica) {
				containsIPs = true
			}
		}
	}

	if len(config.SqlSettings.DataSourceSearchReplicas) > 0 {
		for _, replica := range config.SqlSettings.DataSourceReplicas {
			if ipRegexp.MatchString(replica) {
				containsIPs = true
			}
		}
	}

	if !containsIPs {
		result.Result = check.Result.Pass
		result.Status = Pass

	}

	return result
}

//
//
// EMAIL SETTINGS
//
//

func p002(config model.Config, checks map[string]Check) CheckResult {
	check := checks["p002"]

	result := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Description: check.Description,
		Result:      fmt.Sprintf(check.Result.Fail, *config.EmailSettings.PushNotificationContents),
		Status:      Warn,
		Severity:    CheckSeverity(check.Severity),
	}
	if *config.EmailSettings.PushNotificationContents == "id_loaded" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

//
//
// ELASTICSEARCH SETTINGS
//
//

func h002(config model.Config, checks map[string]Check) CheckResult {
	check := checks["h002"]

	result := CheckResult{
		Name:        check.Name,
		Type:        check.Type,
		Description: check.Description,
		Result:      check.Result.Fail,
		Status:      Warn,
		Severity:    CheckSeverity(check.Severity),
	}

	if *config.ElasticsearchSettings.EnableIndexing {
		if *config.ElasticsearchSettings.LiveIndexingBatchSize > 1 {
			result.Result = fmt.Sprintf(check.Result.Pass, *config.ElasticsearchSettings.LiveIndexingBatchSize)
			result.Status = Pass
		}
		return result
	}

	result.Status = Ignore
	result.Result = check.Result.Ignore

	return result
}
