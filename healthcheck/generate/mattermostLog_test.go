package processpacket

import (
	"testing"
)

func TestH003(t *testing.T) {
	p, checkStatus := setupTest(t, "mattermostLog")

	testCases := []struct {
		name           string
		logs           []byte
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h003 - logs contain 'context deadline exceeded'",
			logs:           []byte("context deadline exceeded"),
			expectedStatus: Fail,
			expectedResult: "Found",
		},
		{
			name:           "h003 - logs do not contain 'context deadline exceeded'",
			logs:           []byte(""),
			expectedStatus: Pass,
			expectedResult: "Not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Logs = tc.logs
			checkStatus(t, p.h003, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH004(t *testing.T) {
	p, checkStatus := setupTest(t, "mattermostLog")

	testCases := []struct {
		name           string
		logs           []byte
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h004 - logs contain 'i/o timeout'",
			logs:           []byte("i/o timeout"),
			expectedStatus: Fail,
			expectedResult: "Found",
		},
		{
			name:           "h004 - logs do not contain 'i/o timeout'",
			logs:           []byte(""),
			expectedStatus: Pass,
			expectedResult: "Not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Logs = tc.logs
			checkStatus(t, p.h004, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}

func TestH005(t *testing.T) {
	p, checkStatus := setupTest(t, "mattermostLog")

	testCases := []struct {
		name           string
		logs           []byte
		expectedStatus CheckStatus
		expectedResult string
	}{
		{
			name:           "h005 - logs contain 'Error while creating session for user access token'",
			logs:           []byte("Error while creating session for user access token"),
			expectedStatus: Fail,
			expectedResult: "Found",
		},
		{
			name:           "h005 - logs do not contain 'Error while creating session for user access token'",
			logs:           []byte(""),
			expectedStatus: Pass,
			expectedResult: "Not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p.packet.Logs = tc.logs
			checkStatus(t, p.h005, nil, tc.expectedStatus, tc.expectedResult)
		})
	}
}
