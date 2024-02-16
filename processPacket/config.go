package processpacket

import (
	"fmt"
	"regexp"

	md "github.com/go-spectest/markdown"
	"github.com/mattermost/mattermost/server/public/model"
)

type CheckFunc func(config model.Config) CheckResult

func configChecks(config model.Config, results *md.Markdown) {

	checks := []CheckFunc{siteURL, extendSessionLengthWithActivity, idNotifications, elasticSearchLiveIndexing, enableLinkPreviews, ipsInSqlConfig}
	testResults := []CheckResult{}

	for _, check := range checks {
		result := check(config)
		testResults = append(testResults, result)
	}

	resultsToArray := [][]string{}

	for _, result := range testResults {
		resultsToArray = append(resultsToArray, []string{result.Name, string(result.Type), result.Status, result.Result, result.Description})
	}

	fmt.Println(resultsToArray)
	results.
		H2("Configuration Checks").
		CustomTable(md.TableSet{
			Header: []string{"Name", "Type", "Status", "Result", "Description"},
			Rows:   resultsToArray,
		}, md.TableOptions{
			AutoWrapText: false,
		})
}

//
//
// SERVICE SETTINGS
//
//

func siteURL(config model.Config) CheckResult {
	result := CheckResult{
		Name:        "Site URL",
		Result:      "Site URL is not set",
		Type:        Health,
		Description: "The siteURL is required by many functions of Mattermost. With it not set some features may not work as expected. [documentation](https://docs.mattermost.com/configure/web-server-configuration-settings.html#site-url)",
		Status:      "🔴",
	}
	if *config.ServiceSettings.SiteURL != "" {
		result.Result = "Site URL is set"
		result.Status = "🟢"
	}
	return result
}

func enableLinkPreviews(config model.Config) CheckResult {
	result := CheckResult{
		Name:        "Link Previews",
		Type:        Adoption,
		Description: "Link Previews are a feature that allows for a preview of a link to be displayed in the Mattermost client to improve end user experience. [documentation](https://docs.mattermost.com/configure/site-configuration-settings.html#posts-enablemessagelinkpreviews)",
		Status:      "🟡",
		Result:      "Link Previews are not enabled",
	}

	if *config.ServiceSettings.EnableLinkPreviews {
		result.Result = "Link Previews are enabled"
		result.Status = "🟢"
	}

	return result
}

func extendSessionLengthWithActivity(config model.Config) CheckResult {
	result := CheckResult{
		Name:        "Extend Session Length with Activity",
		Type:        Adoption,
		Description: "For improved end-user login session lifecycle, consider enableing `ExtendSessionLengthWithActivity` Verify with Enterprise policies if this is compatible. [documenation](https://docs.mattermost.com/configure/session-lengths-configuration-settings.html#extend-session-length-with-activity)",
		Status:      "🟡",
	}

	if *config.ServiceSettings.ExtendSessionLengthWithActivity {
		result.Result = "Session Length is extended with activity"
		result.Status = "🟢"
	}

	return result
}

//
//
// SQL SETTINGS
//
//

func ipsInSqlConfig(config model.Config) CheckResult {
	result := CheckResult{
		Name:        "IPs in SQL Data Sources",
		Type:        Proactive,
		Description: "Using IP addresses in your SQL data sources can cause issues with failovers in the event of a database failure. [documentation](https://docs.mattermost.com/configure/environment-configuration-settings.html#database-datasource)",
		Status:      "🟡",
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

	if containsIPs {
		result.Result = "Data sources contain IP addresses"
		result.Status = "🔴"
	}

	if !containsIPs {
		result.Result = "Data sources do not contain IP addresses"
		result.Status = "🟢"

	}

	return result
}

//
//
// EMAIL SETTINGS
//
//

func idNotifications(config model.Config) CheckResult {
	result := CheckResult{
		Name:        "ID-Only Notifications",
		Type:        Proactive,
		Description: "Setting notifications to ID-Only keeps data off Google / Apple servers and in turn your server is more secure. [documentation](https://docs.mattermost.com/configure/site-configuration-settings.html#push-notification-contents)",
		Status:      "🟡",
	}
	if *config.EmailSettings.EmailNotificationContentsType == "id_loaded" {
		result.Result = "ID Notifications is set to `id_loaded`"
		result.Status = "🟢"
	} else {

		result.Result = "Notification contents are set to `" + *config.EmailSettings.EmailNotificationContentsType + "`"
		result.Status = "🟡"
	}

	return result
}

//
//
// ELASTICSEARCH SETTINGS
//
//

func elasticSearchLiveIndexing(config model.Config) CheckResult {
	result := CheckResult{
		Name:        "ElasticSearch Live Indexing",
		Type:        Health,
		Description: "Live index batch size controls how often Elasticsearch is indexed in real time. On highly active servers this needs to be increased to prevent an Elasticsearch crash. [documentation](https://docs.mattermost.com/configure/environment-configuration-settings.html#live-indexing-batch-size)",
		Status:      "🔴",
	}

	if *config.ElasticsearchSettings.EnableIndexing {
		if *config.ElasticsearchSettings.LiveIndexingBatchSize > 1 {
			result.Result = "Live Indexing has been modified to a value greater than default of 1 - ` " + fmt.Sprint(*config.ElasticsearchSettings.LiveIndexingBatchSize) + " `"
			result.Status = "🟢"
		} else {
			result.Result = "ElasticSearch Live Indexing has not been modified"
			result.Status = "🔴"
		}
		return result
	}

	result.Status = "-"
	result.Result = "Elasticsearch is not enabled"

	return result
}
