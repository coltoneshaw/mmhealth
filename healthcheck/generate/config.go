package processpacket

import (
	"fmt"
	"regexp"

	"github.com/mattermost/mattermost/server/public/model"
)

type ConfigCheckFunc func(config model.Config, checks map[string]Check) CheckResult

func (p *ProcessPacket) configChecks(config model.Config) (results []CheckResult) {

	checks := map[string]ConfigCheckFunc{
		"h001": p.h001,
		"h002": p.h002,
		"p001": p.p001,
		"p002": p.p002,
		"a001": p.a001,
		"a002": p.a002,
	}
	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(config, p.Checks.Config)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

//
//
// SERVICE SETTINGS
//
//

func (p *ProcessPacket) h001(config model.Config, checks map[string]Check) CheckResult {
	check, result := initCheckResult("h001", checks, Fail)
	if *config.ServiceSettings.SiteURL != "" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}
	return result
}

func (p *ProcessPacket) a001(config model.Config, checks map[string]Check) CheckResult {
	check, result := initCheckResult("a001", checks, Fail)

	if *config.ServiceSettings.EnableLinkPreviews {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

func (p *ProcessPacket) a002(config model.Config, checks map[string]Check) CheckResult {
	check, result := initCheckResult("a002", checks, Fail)

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

func (p *ProcessPacket) p001(config model.Config, checks map[string]Check) CheckResult {
	check, result := initCheckResult("p001", checks, Fail)

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

func (p *ProcessPacket) p002(config model.Config, checks map[string]Check) CheckResult {
	check, result := initCheckResult("p002", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, *config.EmailSettings.PushNotificationContents)

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

func (p *ProcessPacket) h002(config model.Config, checks map[string]Check) CheckResult {
	check, result := initCheckResult("h002", checks, Fail)
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
