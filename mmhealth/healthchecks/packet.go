package healthchecks

import (
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
)

func (p *ProcessPacket) packetChecks() (results []CheckResult) {

	checks := map[string]CheckFunc{
		"h012": p.h012,
	}

	testResults := []CheckResult{}

	for id, check := range checks {
		result := check(p.Checks.Packet)
		result.ID = id
		testResults = append(testResults, result)
	}

	return p.sortResults(testResults)
}

// checks to see if any of the ldap sync jobs have failed and if LDAP is enabled. If so we fail the job.
func (p *ProcessPacket) h012(checks map[string]types.Check) CheckResult {
	// check defaults to pass here because we are looking for the failure message
	check, result := initCheckResult("h012", checks, Pass)

	// check if LDAP is enabled in the config
	if !*p.packet.Config.LdapSettings.Enable {
		result.Result = check.Result.Ignore
		result.Status = Ignore
		return result
	}

	// check if the ldap_sync_jobs for any status that's not success
	for _, job := range p.packet.Packet.LdapSyncJobs {
		if job.Status != "success" {
			result.Result = check.Result.Fail
			result.Status = Fail
			return result
		}
	}

	return result
}
