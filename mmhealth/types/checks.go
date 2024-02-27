package types

import "github.com/mattermost/mattermost/server/public/model"

type CheckType string

const (
	Proactive CheckType = "proactive"
	Health    CheckType = "health"
	Adoption  CheckType = "adoption"
)

type CheckStatus string

const (
	Fail   CheckStatus = "fail"
	Pass   CheckStatus = "pass"
	Warn   CheckStatus = "warn"
	Ignore CheckStatus = "ignore"
)

type CheckSeverity string

const (
	Urgent CheckSeverity = "urgent"
	High   CheckSeverity = "high"
	Medium CheckSeverity = "medium"
	Low    CheckSeverity = "low"
)

type Result struct {
	Pass   string `yaml:"pass"`
	Fail   string `yaml:"fail"`
	Ignore string `yaml:"ignore"`
}

type Check struct {
	Name        string        `yaml:"name"`
	Result      Result        `yaml:"result"`
	Description string        `yaml:"description"`
	Severity    CheckSeverity `yaml:"severity"`
	Type        CheckType     `yaml:"type"`
}

type ChecksFile struct {
	Config          map[string]Check `yaml:"config"`
	Packet          map[string]Check `yaml:"packet"`
	Environment     map[string]Check `yaml:"environment"`
	MattermostLog   map[string]Check `yaml:"mattermostLog"`
	NotificationLog map[string]Check `yaml:"notificationLog"`
	Plugins         map[string]Check `yaml:"plugins"`
}

type PacketData struct {
	Logs             []MattermostLogEntry
	NotificationLogs []byte
	Config           model.Config
	Plugins          model.PluginsResponse
	Packet           model.SupportPacket
}

type ConfigFile struct {
	Versions struct {
		Supported []string `yaml:"supported"`
		ESR       string   `yaml:"esr"`
	} `yaml:"versions"`
}
