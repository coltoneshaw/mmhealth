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
