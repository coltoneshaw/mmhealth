package healthchecks

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
)

func (p *ProcessPacket) pluginChecks() (results []PluginResults) {

	plugins := p.packet.Plugins
	pluginResults := []PluginResults{}

	for _, plugin := range plugins.Active {
		parsedPlugin := getPluginResults(plugin, p.Config, true)
		pluginResults = append(pluginResults, parsedPlugin)
	}

	for _, plugin := range plugins.Inactive {
		parsedPlugin := getPluginResults(plugin, p.Config, false)
		pluginResults = append(pluginResults, parsedPlugin)
	}

	return pluginResults

}

func getPluginResults(plugin *model.PluginInfo, config types.ConfigFile, isActive bool) PluginResults {
	parsedPlugin := PluginResults{
		PluginID:         plugin.Id,
		PluginName:       plugin.Name,
		PluginURL:        findAPluginURL(plugin, config),
		Active:           isActive,
		LatestVersion:    findLatestVersion(plugin.Id, config),
		InstalledVersion: plugin.Version,
	}

	isUpdated, err := returnIsUpdated(parsedPlugin.LatestVersion, parsedPlugin.InstalledVersion)

	if err != nil {
		fmt.Printf("Error parsing version for plugin %s: %s\n", plugin.Id, err)
	}

	parsedPlugin.IsUpdated = isUpdated
	return parsedPlugin
}
func returnIsUpdated(latest string, installed string) (bool, error) {

	installedVersion, err := semver.NewVersion(installed)
	if err != nil {
		return false, err
	}

	latestSemVersion, err := semver.NewVersion(latest)
	if err != nil {
		return false, err
	}

	isUpdated := installedVersion.Compare(latestSemVersion) >= 0

	return isUpdated, nil

}

func findAPluginURL(plugin *model.PluginInfo, config types.ConfigFile) string {
	for ID, p := range config.Plugins {
		if ID == plugin.Id {
			return p.Repo
		}
	}

	if plugin.HomepageURL != "" {
		return plugin.HomepageURL
	}
	if plugin.SupportURL != "" {
		return plugin.SupportURL
	}
	if plugin.ReleaseNotesURL != "" {
		return plugin.ReleaseNotesURL
	}

	return ""
}

func findLatestVersion(pluginID string, config types.ConfigFile) string {

	for ID, plugin := range config.Plugins {
		if ID == pluginID {
			return plugin.Latest
		}
	}

	return "-"
}
