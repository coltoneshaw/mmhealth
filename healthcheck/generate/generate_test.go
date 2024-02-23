package processpacket

import (
	"os"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"gopkg.in/yaml.v3"
)

func testReadChecksFile(t *testing.T) ChecksFile {
	// Save the current working directory
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current dir. err: %v", err)
	}

	// Change the working directory to the directory of this test file
	err = os.Chdir("../../")
	if err != nil {
		t.Fatalf("error changing current dir. err: %v", err)
	}

	// Make sure to change it back when the test finishes
	defer func() {
		err := os.Chdir(oldwd)
		if err != nil {
			t.Fatalf("error changing current back to oldwd: %v", err)
		}
	}()

	// Now you can call your function under test
	var checks ChecksFile
	data, err := os.ReadFile("checks.yaml")
	if err != nil {
		t.Fatalf("error reading checks file: %v", err)
	}

	err = yaml.Unmarshal(data, &checks)
	if err != nil {
		t.Fatalf("error unmarshaling checks file: %v", err)
	}

	return checks
}

func testReadConfigFile(t *testing.T) ConfigFile {
	// Save the current working directory
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting current dir. err: %v", err)
	}

	// Change the working directory to the directory of this test file
	err = os.Chdir("../../")
	if err != nil {
		t.Fatalf("error changing current dir. err: %v", err)
	}

	// Make sure to change it back when the test finishes
	defer func() {
		err := os.Chdir(oldwd)
		if err != nil {
			t.Fatalf("error changing current back to oldwd: %v", err)
		}
	}()

	// Now you can call your function under test
	var config ConfigFile
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		t.Fatalf("error reading config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		t.Fatalf("error unmarshaling config file: %v", err)
	}

	return config
}

func mockProcessPacket(t *testing.T) (*ProcessPacket, error) {
	p := &ProcessPacket{}

	p.packet.Config = model.Config{}
	p.packet.Config.SetDefaults()

	p.Checks = testReadChecksFile(t)
	p.Config = testReadConfigFile(t)

	return p, nil
}

func setupTest(t *testing.T, checkType string) (
	*ProcessPacket,
	func(t *testing.T, testFunc func(checks map[string]Check) CheckResult, expectedStatus CheckStatus, expectedResult string)) {
	p, err := mockProcessPacket(t)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	checks := map[string]Check{}

	if checkType == "" || (checkType != "config" && checkType != "environment" && checkType != "mattermostLog") {
		t.Fatalf("checkType is incorrect")
	}

	switch checkType {
	case "config":
		checks = p.Checks.Config
	case "environment":
		checks = p.Checks.Environment
	case "mattermostLog":
		checks = p.Checks.MattermostLog

	}

	checkStatus := func(t *testing.T, testFunc func(checks map[string]Check) CheckResult, expectedStatus CheckStatus, expectedResult string) {
		test := testFunc(checks)

		if test.Status != expectedStatus {
			t.Errorf("Expected status '%v', got '%v'.", expectedStatus, test.Status)
		}

		if test.Result != expectedResult {
			t.Errorf("Expected result '%v', got '%v'.", expectedResult, test.Result)
		}
	}

	return p, checkStatus
}
