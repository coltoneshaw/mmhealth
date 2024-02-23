package processpacket

import (
	"testing"
)

func TestH001(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		siteURL        string
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h001 - SiteURL not set",
			siteURL:        "",
			expectedStatus: Fail,
			expectedResult: "Not set",
		},
		{
			name:           "h001 - SiteURL set",
			siteURL:        "http://localhost",
			expectedStatus: Pass,
			expectedResult: "Set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.SiteURL = &tc.siteURL
			checkStatus(t, p.h001, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}
func TestA001(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		linkPreviews   bool
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "a001 - Enable link previews is false",
			linkPreviews:   false,
			expectedStatus: Warn,
			expectedResult: "Not enabled",
		},
		{
			name:           "a001 - Enable link previews is true",
			linkPreviews:   true,
			expectedStatus: Pass,
			expectedResult: "Enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.EnableLinkPreviews = &tc.linkPreviews
			checkStatus(t, p.a001, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestA002(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		sessionLength  bool
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "a002 - Extend session length with activity is false",
			sessionLength:  false,
			expectedStatus: Warn,
			expectedResult: "Not enabled",
		},
		{
			name:           "a002 - Extend session length with activity is true",
			sessionLength:  true,
			expectedStatus: Pass,
			expectedResult: "Enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.ExtendSessionLengthWithActivity = &tc.sessionLength
			checkStatus(t, p.a002, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}
func TestP002(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name             string
		notificationType string
		expectedStatus   CheckStatus
		expectedResult   string
	}{
		{
			name:             "p002 - not using ID notifications",
			notificationType: "full",
			expectedStatus:   Warn,
			expectedResult:   "Set to `full`",
		},
		{
			name:             "p002 - using ID notifications",
			notificationType: "id_loaded",
			expectedStatus:   Pass,
			expectedResult:   "Set to `id_loaded`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.EmailSettings.PushNotificationContents = &tc.notificationType
			checkStatus(t, p.p002, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH002(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		enableIndexing bool
		liveIndexing   int
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h002 - ES not in use",
			enableIndexing: false,
			liveIndexing:   10,
			expectedStatus: Ignore,
			expectedResult: "Elasticsearch disabled",
		},
		{
			name:           "h002 - ES in use and live indexing is set to default",
			enableIndexing: true,
			liveIndexing:   1,
			expectedStatus: Fail,
			expectedResult: "Uses default value",
		},
		{
			name:           "h002 - ES in use and live indexing is configured",
			enableIndexing: true,
			liveIndexing:   10,
			expectedStatus: Pass,
			expectedResult: "Modified to a value greater than default of 1 - `10`",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.LiveIndexingBatchSize = &tc.liveIndexing
			checkStatus(t, p.h002, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH010(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name                  string
		enableIndexing        bool
		enableSearching       bool
		enableAutocomplete    bool
		disableDatabaseSearch bool
		expectedStatus        CheckStatus
		expectedResult        string
	}{
		{
			name:                  "h010 - using Database search",
			enableIndexing:        false,
			enableSearching:       false,
			enableAutocomplete:    false,
			disableDatabaseSearch: false,
			expectedStatus:        Pass,
			expectedResult:        "Database",
		},
		{
			name:                  "h010 - using Elasticsearch",
			enableIndexing:        true,
			enableSearching:       true,
			enableAutocomplete:    true,
			disableDatabaseSearch: false,
			expectedStatus:        Pass,
			expectedResult:        "Elasticsearch",
		},
		{
			name:                  "h010 - No search enabled",
			enableIndexing:        false,
			enableSearching:       false,
			enableAutocomplete:    false,
			disableDatabaseSearch: true,
			expectedStatus:        Fail,
			expectedResult:        "No search enabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.EnableSearching = &tc.enableSearching
			p.packet.Config.ElasticsearchSettings.EnableAutocomplete = &tc.enableAutocomplete
			p.packet.Config.SqlSettings.DisableDatabaseSearch = &tc.disableDatabaseSearch
			checkStatus(t, p.h010, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestP003(t *testing.T) {
	p, checkStatus := setupTest(t, "config")

	testCases := []struct {
		name           string
		IdAttribute    string
		ldapEnabled    bool
		emailAttribute string
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "p003 - not using LDAP",
			IdAttribute:    "",
			emailAttribute: "",
			ldapEnabled:    false,
			expectedStatus: Ignore,
			expectedResult: "LDAP disabled",
		},
		{
			name:           "p003 - LDAP enabled and ID attribute set",
			IdAttribute:    "uniqueID",
			emailAttribute: "email",
			ldapEnabled:    true,
			expectedStatus: Pass,
			expectedResult: "ID attribute set to uniqueID.",
		},
		{
			name:           "p003 - LDAP enabled and ID attribute set the same",
			IdAttribute:    "myEmail",
			emailAttribute: "myEmail",
			ldapEnabled:    true,
			expectedStatus: Warn,
			expectedResult: "Using email",
		},
		{
			name:           "p003 - LDAP enabled no ID set",
			IdAttribute:    "",
			emailAttribute: "",
			ldapEnabled:    true,
			expectedStatus: Warn,
			expectedResult: "Using email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.LdapSettings.Enable = &tc.ldapEnabled
			p.packet.Config.LdapSettings.IdAttribute = &tc.IdAttribute
			p.packet.Config.LdapSettings.EmailAttribute = &tc.emailAttribute

			checkStatus(t, p.p003, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}
