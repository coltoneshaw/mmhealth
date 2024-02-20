package processpacket

import (
	"testing"
)

func TestH001(t *testing.T) {
	p, checkStatus := setupTest(t)

	testCases := []struct {
		name           string
		siteURL        string
		expectedStatus CheckStatus
	}{
		{
			name:           "h001 - SiteURL not set",
			siteURL:        "",
			expectedStatus: Fail,
		},
		{
			name:           "h001 - SiteURL set",
			siteURL:        "http://localhost",
			expectedStatus: Pass,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.SiteURL = &tc.siteURL
			checkStatus(t, p.h001, nil, tc.expectedStatus)
		})
	}
}
func TestA001(t *testing.T) {
	p, checkStatus := setupTest(t)

	testCases := []struct {
		name           string
		linkPreviews   bool
		expectedStatus CheckStatus
	}{
		{
			name:           "a001 - Enable link previews is false",
			linkPreviews:   false,
			expectedStatus: Warn,
		},
		{
			name:           "a001 - Enable link previews is true",
			linkPreviews:   true,
			expectedStatus: Pass,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.EnableLinkPreviews = &tc.linkPreviews
			checkStatus(t, p.a001, nil, tc.expectedStatus)
		})
	}
}

func TestA002(t *testing.T) {
	p, checkStatus := setupTest(t)

	testCases := []struct {
		name           string
		sessionLength  bool
		expectedStatus CheckStatus
	}{
		{
			name:           "a002 - Extend session length with activity is false",
			sessionLength:  false,
			expectedStatus: Warn,
		},
		{
			name:           "a002 - Extend session length with activity is true",
			sessionLength:  true,
			expectedStatus: Pass,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ServiceSettings.ExtendSessionLengthWithActivity = &tc.sessionLength
			checkStatus(t, p.a002, nil, tc.expectedStatus)
		})
	}
}
func TestP002(t *testing.T) {
	p, checkStatus := setupTest(t)

	testCases := []struct {
		name             string
		notificationType string
		expectedStatus   CheckStatus
	}{
		{
			name:             "p002 - not using ID notifications",
			notificationType: "full",
			expectedStatus:   Warn,
		},
		{
			name:             "p002 - using ID notifications",
			notificationType: "id_loaded",
			expectedStatus:   Pass,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.EmailSettings.PushNotificationContents = &tc.notificationType
			checkStatus(t, p.p002, nil, tc.expectedStatus)
		})
	}
}

func TestH002(t *testing.T) {
	p, checkStatus := setupTest(t)

	testCases := []struct {
		name           string
		enableIndexing bool
		liveIndexing   int
		expectedStatus CheckStatus
	}{
		{
			name:           "h002 - ES not in use",
			enableIndexing: false,
			liveIndexing:   10,
			expectedStatus: Ignore,
		},
		{
			name:           "h002 - ES in use and live indexing is set to default",
			enableIndexing: true,
			liveIndexing:   1,
			expectedStatus: Fail,
		},
		{
			name:           "h002 - ES in use and live indexing is configured",
			enableIndexing: true,
			liveIndexing:   10,
			expectedStatus: Pass,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.LiveIndexingBatchSize = &tc.liveIndexing
			checkStatus(t, p.h002, nil, tc.expectedStatus)
		})
	}
}

func TestH010(t *testing.T) {
	p, checkStatus := setupTest(t)

	testCases := []struct {
		name                  string
		enableIndexing        bool
		enableSearching       bool
		enableAutocomplete    bool
		disableDatabaseSearch bool
		expectedStatus        CheckStatus
	}{
		{
			name:                  "h010 - using Database search",
			enableIndexing:        false,
			enableSearching:       false,
			enableAutocomplete:    false,
			disableDatabaseSearch: false,
			expectedStatus:        Pass,
		},
		{
			name:                  "h010 - using Elasticsearch",
			enableIndexing:        true,
			enableSearching:       true,
			enableAutocomplete:    true,
			disableDatabaseSearch: false,
			expectedStatus:        Pass,
		},
		{
			name:                  "h010 - No search enabled",
			enableIndexing:        false,
			enableSearching:       false,
			enableAutocomplete:    false,
			disableDatabaseSearch: true,
			expectedStatus:        Fail,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Config.ElasticsearchSettings.EnableIndexing = &tc.enableIndexing
			p.packet.Config.ElasticsearchSettings.EnableSearching = &tc.enableSearching
			p.packet.Config.ElasticsearchSettings.EnableAutocomplete = &tc.enableAutocomplete
			p.packet.Config.SqlSettings.DisableDatabaseSearch = &tc.disableDatabaseSearch
			checkStatus(t, p.h010, nil, tc.expectedStatus)
		})
	}
}
