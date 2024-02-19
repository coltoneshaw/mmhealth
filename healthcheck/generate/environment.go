package processpacket

import (
	"fmt"

	"github.com/Masterminds/semver"
)

type PacketCheckFunc func(checks map[string]Check) CheckResult

func (p *ProcessPacket) environmentChecks() (results []CheckResult) {

	checks := map[string]PacketCheckFunc{
		"h006": p.h006,
		"h007": p.h007,
		"h008": p.h008,
		"h009": p.h009,
	}

	for id, check := range checks {
		result := check(p.Checks.Environment)
		result.ID = id
		results = append(results, result)
	}

	return p.sortResults(results)
}

// Server Version check
func (p *ProcessPacket) h006(checks map[string]Check) CheckResult {

	check, result := initCheckResult("h006", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Packet.ServerVersion)

	serverVersion, err := semver.NewVersion(p.packet.Packet.ServerVersion)
	if err != nil {
		result.Status = Fail
		result.Result = fmt.Sprintf("Error parsing server version: %s", err)
		return result
	}

	currentVersion, _ := semver.NewVersion(p.Config.Versions.Current)
	currentESR, _ := semver.NewVersion(p.Config.Versions.Esr)

	diff := currentVersion.Minor() - serverVersion.Minor()

	// checking if the server version is within two of the current release
	if diff >= 0 && diff <= 2 {
		result.Result = fmt.Sprintf(check.Result.Pass, p.packet.Packet.ServerVersion)
		result.Status = Pass
	}

	if serverVersion.Major() == currentESR.Major() && serverVersion.Minor() == currentESR.Minor() {
		result.Result = fmt.Sprintf(check.Result.Pass, p.packet.Packet.ServerVersion)
		result.Status = Pass
	}

	return result
}

func (p *ProcessPacket) h007(checks map[string]Check) CheckResult {
	check, result := initCheckResult("h007", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Packet.DatabaseType)

	if p.packet.Packet.DatabaseType == "postgres" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

func (p *ProcessPacket) h008(checks map[string]Check) CheckResult {
	check, result := initCheckResult("h008", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Packet.ServerOS)

	if p.packet.Packet.ServerOS == "linux" {
		result.Result = check.Result.Pass
		result.Status = Pass
	}

	return result
}

func (p *ProcessPacket) h009(checks map[string]Check) CheckResult {
	check, result := initCheckResult("h009", checks, Fail)

	result.Result = fmt.Sprintf(check.Result.Fail, p.packet.Packet.TotalPosts, *p.packet.Config.ElasticsearchSettings.EnableIndexing)

	if p.packet.Packet.TotalPosts < 2500000 || *p.packet.Config.ElasticsearchSettings.EnableIndexing {
		result.Result = fmt.Sprintf(check.Result.Pass, p.packet.Packet.TotalPosts, *p.packet.Config.ElasticsearchSettings.EnableIndexing)
		result.Status = Pass
	}

	return result
}
