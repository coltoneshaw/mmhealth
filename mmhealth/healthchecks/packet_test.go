package healthchecks

import (
	"testing"

	"github.com/coltoneshaw/mmhealth/mmhealth/types"
	"github.com/mattermost/mattermost/server/public/model"
)

func TestH012(t *testing.T) {
	p, checkStatus := setupTest(t, "packet")

	testCases := []struct {
		name           string
		ldapEnabled    bool
		expectedStatus types.CheckStatus
		expectedResult string
		jobs           []*model.Job
	}{
		{
			name:           "h012 - LDAP is not enabled",
			ldapEnabled:    false,
			expectedStatus: Ignore,
			expectedResult: "LDAP is disabled",
			jobs:           []*model.Job{},
		},
		{
			name:           "h012 - LDAP enabled with passed job",
			ldapEnabled:    true,
			expectedStatus: Pass,
			expectedResult: "LDAP jobs succeeded",
			jobs: []*model.Job{
				{
					Status: "success",
				},
			},
		},
		{
			name:           "h012 - LDAP enabled with failed job",
			ldapEnabled:    true,
			expectedStatus: Fail,
			expectedResult: "LDAP jobs failed",
			jobs: []*model.Job{
				{
					Status: "failed",
				},
			},
		},
		{
			name:           "h012 - LDAP disabled with failed job",
			ldapEnabled:    false,
			expectedStatus: Ignore,
			expectedResult: "LDAP is disabled",
			jobs: []*model.Job{
				{
					Status: "failed",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.LdapSettings.Enable = &tc.ldapEnabled
			p.packet.Packet.LdapSyncJobs = tc.jobs
			checkStatus(t, p.h012, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH013(t *testing.T) {
	p, checkStatus := setupTest(t, "packet")

	testCases := []struct {
		name                 string
		messageExportEnabled bool
		expectedStatus       types.CheckStatus
		expectedResult       string
		jobs                 []*model.Job
	}{
		{
			name:                 "h013 - message export is not enabled",
			messageExportEnabled: false,
			expectedStatus:       Ignore,
			expectedResult:       "Message export is disabled",
			jobs:                 []*model.Job{},
		},
		{
			name:                 "h013 - message export enabled with passed job",
			messageExportEnabled: true,
			expectedStatus:       Pass,
			expectedResult:       "Message export jobs succeeded",
			jobs: []*model.Job{
				{
					Status: "success",
				},
			},
		},
		{
			name:                 "h013 - message export enabled with failed job",
			messageExportEnabled: true,
			expectedStatus:       Fail,
			expectedResult:       "Message export jobs failed",
			jobs: []*model.Job{
				{
					Status: "failed",
				},
			},
		},
		{
			name:                 "h013 - message export disabled with failed job",
			messageExportEnabled: false,
			expectedStatus:       Ignore,
			expectedResult:       "Message export is disabled",
			jobs: []*model.Job{
				{
					Status: "failed",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.MessageExportSettings.EnableExport = &tc.messageExportEnabled
			p.packet.Packet.MessageExportJobs = tc.jobs
			checkStatus(t, p.h013, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH014(t *testing.T) {
	p, checkStatus := setupTest(t, "packet")

	testCases := []struct {
		name           string
		expectedStatus types.CheckStatus
		expectedResult string
		jobs           []*model.Job
	}{
		{
			name:           "h014 - No migration jobs found",
			expectedStatus: Ignore,
			expectedResult: "No migration jobs found",
			jobs:           []*model.Job{},
		},
		{
			name:           "h014 - Migration jobs succeeded",
			expectedStatus: Pass,
			expectedResult: "Migration jobs succeeded",
			jobs: []*model.Job{
				{
					Status: "success",
				},
			},
		},
		{
			name:           "h014 - Migration jobs failed",
			expectedStatus: Fail,
			expectedResult: "Migration jobs failed",
			jobs: []*model.Job{
				{
					Status: "failed",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Packet.MigrationJobs = tc.jobs
			checkStatus(t, p.h014, tc.expectedStatus, tc.expectedResult)
		})
	}
}
