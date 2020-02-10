package urn

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURN(t *testing.T) {
	testCases := []struct {
		name string
		URN  string

		expectedURN   *URN
		expectedError error
	}{
		{
			name: "valid URN",
			URN:  "urn:newtonworld228:lol%AC_45:rRR",
			expectedURN: &URN{
				[]byte("newtonworld228"),
				[]byte("lol%AC_45:rRR"),
			},
			expectedError: nil,
		},
		{
			name:          "empty URN",
			URN:           "",
			expectedURN:   nil,
			expectedError: errors.New("invalid URN format, should be urn:<nid>:<nss>"),
		},
		{
			name:          "wrong prefix in URN",
			URN:           "irn:nid456:nss_$:kek",
			expectedURN:   nil,
			expectedError: errors.New("URN 'irn:nid456:nss_$:kek' must have prefix - urn"),
		},
		{
			name:          "wrong NID",
			URN:           "urn:$sdf:nss_$:kek",
			expectedURN:   nil,
			expectedError: errors.New("NID $sdf doesn't satisfy pattern: ^[a-zA-Z0-9]{1}[a-zA-Z0-9\\-]{1,31}$"),
		},
		{
			name:        "wrong NSS",
			URN:         "urn:newtonworld228:?kek:lol",
			expectedURN: nil,
			expectedError: errors.New(
				"NSS ?kek:lol doesn't satisfy the regexp rule: ^(?:%[0-9A-Fa-f]{2}|[a-zA-Z0-9\\-+(),.:=@;$_!*'])+$"),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				u, err := Parse(testCase.URN)
				assert.Equal(t, testCase.expectedURN, u)
				assert.Equal(t, testCase.expectedError, err)
			},
		)
	}
}

func BenchmarkParseURN(b *testing.B) {
	rawURN := []byte("urn:newtonworld228:myNSSwith_%23HEXvalue:kek")

	for i := 0; i < b.N; i++ {
		_, err := parseRawURN(rawURN)
		if err != nil {
			fmt.Printf(
				"error: %s\n", err,
			)
		}
	}
}

func BenchmarkURNToStringTest(b *testing.B) {
	u, _ := New("newtonworld228", "myNSSwith_%23HEXvalue:kek")

	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}

func BenchmarkURNToJSONTest(b *testing.B) {
	u, _ := New("newtonworld228", "myNSSwith_%23HEXvalue:kek")

	for i := 0; i < b.N; i++ {
		_, err := u.MarshalJSON()
		if err != nil {
			fmt.Printf(
				"error: %s\n", err,
			)
		}
	}
}

func BenchmarkJSONMarshalTest(b *testing.B) {
	u, _ := New("newtonworld228", "myNSSwith_%23HEXvalue:kek")

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(u)
		if err != nil {
			fmt.Printf(
				"error: %s\n", err,
			)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	someURN := &URN{[]byte("newtonworld228"), []byte("myNSSwith_%23HEXvalue:kek")}
	someStruct := struct {
		Urn *URN `json:"myurn"`
	}{Urn: someURN}

	expectedResult := []byte(`{"myurn":"urn:newtonworld228:myNSSwith_%23HEXvalue:kek"}`)

	marshalled, err := json.Marshal(someStruct)

	assert.Equal(t, err, nil)
	assert.Equal(t, expectedResult, marshalled)
}

func TestNewURN(t *testing.T) {
	testCases := []struct {
		name string
		nid  string
		nss  string

		expectedURN       *URN
		expectedStringURN string
		expectedError     error
	}{
		{
			name: "valid nid and nss",
			nid:  "newtonworld",
			nss:  "user:test_-user",
			expectedURN: &URN{
				nid: []byte("newtonworld"),
				nss: []byte("user:test_-user"),
			},
			expectedStringURN: "urn:newtonworld:user:test_-user",
			expectedError:     nil,
		},
		{
			name: "valid nid and nss with HEX encoded",
			nid:  "newtonworld",
			nss:  "user:test_-user%2C",
			expectedURN: &URN{
				nid: []byte("newtonworld"),
				nss: []byte("user:test_-user%2C"),
			},
			expectedStringURN: "urn:newtonworld:user:test_-user%2C",
			expectedError:     nil,
		},
		{
			name:              "too short nid and valid nss",
			nid:               "n",
			nss:               "user:test_-user",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: length of NID must be more than 2 letters long"),
		},
		{
			name:              "too long nid and valid nss",
			nid:               "nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn",
			nss:               "user:test_-user",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NID must be mot greater than 32 letters long"),
		},
		{
			name:              "reserved nid and valid nss",
			nid:               "urn-nid",
			nss:               "user:test_-user",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NID urn- is reserved"),
		},
		{
			name:              "experimental nid and valid nss",
			nid:               "x-nid",
			nss:               "user:test_-user",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NID x- is experimental"),
		},
		{
			name:              "experimental xy- nid and valid nss",
			nid:               "XY-nid",
			nss:               "user:test_-user",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NID XY-nid mustn't start with: xy-"),
		},
		{
			name:              "wrong characters in nid and valid nss",
			nid:               "_$nid",
			nss:               "user:test_-user",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NID _$nid doesn't satisfy pattern: ^[a-zA-Z0-9]{1}[a-zA-Z0-9\\-]{1,31}$"),
		},
		{
			name:              "valid nid and empty nss",
			nid:               "newtonworld",
			nss:               "",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NSS must be at least one character long"),
		},
		{
			name:              "valid nid and invalid nss",
			nid:               "newtonworld",
			nss:               "%%lol?kek",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: NSS doesn't satisfy the regexp rule: ^(?:%[0-9A-Fa-f]{2}|[a-zA-Z0-9\\-+(),.:=@;$_!*'])+$"),
		},
		{
			name:              "empty nid and nss",
			nid:               "",
			nss:               "",
			expectedURN:       nil,
			expectedStringURN: "",
			expectedError:     errors.New("can't create URN, reason: length of NID must be more than 2 letters long"),
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				u, err := New(
					testCase.nid,
					testCase.nss,
				)

				assert.Equal(t, testCase.expectedURN, u)

				if u != nil {
					assert.Equal(t, testCase.expectedStringURN, u.String())
				}
				assert.Equal(t, testCase.expectedError, err)
			},
		)
	}
}
