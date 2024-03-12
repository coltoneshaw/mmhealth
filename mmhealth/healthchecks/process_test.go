package healthchecks

import (
	"runtime"
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/files"
	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
)

func mockProcessPacket(t *testing.T) (*ProcessPacket, error) {
	p := &ProcessPacket{}

	p.packet.Config = model.Config{}
	p.packet.Config.SetDefaults()

	checks, err := files.ReadChecksFile()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	p.Checks = checks

	config, err := files.ReadConfigFile()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	p.Config = config

	return p, nil
}

func setupTest(t *testing.T, checkType string) (
	*ProcessPacket,
	func(t *testing.T, testFunc func(checks map[string]types.Check) CheckResult, expectedStatus types.CheckStatus, expectedResult string)) {
	p, err := mockProcessPacket(t)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	checks := map[string]types.Check{}

	if checkType == "" || (checkType != "config" && checkType != "environment" && checkType != "mattermostLog" && checkType != "packet") {
		t.Fatalf("checkType is incorrect")
	}

	switch checkType {
	case "config":
		checks = p.Checks.Config
	case "environment":
		checks = p.Checks.Environment
	case "mattermostLog":
		checks = p.Checks.MattermostLog
	case "packet":
		checks = p.Checks.Packet

	}

	checkStatus := func(t *testing.T, testFunc func(checks map[string]types.Check) CheckResult, expectedStatus types.CheckStatus, expectedResult string) {
		test := testFunc(checks)
		_, file, line, _ := runtime.Caller(1)

		if test.Status != expectedStatus {
			t.Errorf("Failed at %s:%d: Expected status '%v', got '%v'.", file, line, expectedStatus, test.Status)
		}

		if test.Result != expectedResult {
			t.Errorf("Failed at %s:%d: Expected result '%v', got '%v'.", file, line, expectedResult, test.Result)
		}
	}

	return p, checkStatus
}
