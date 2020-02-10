package urn

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

var (
	// https://tools.ietf.org/html/rfc3406#section-4.1
	// https://tools.ietf.org/html/rfc3406#section-4.3
	reservedNIDPrefix     = []byte{117, 114, 110, 45} // urn-
	experimentalNIDPrefix = []byte{120, 45}           // x-
	xyNIDPrefix           = []byte{120, 121, 45}      // xy-

	urnPrefix    = []byte{117, 114, 110} //urn
	urnDelimiter = []byte{58}            // :

	// https://tools.ietf.org/html/rfc2141
	nidRegexp = regexp.MustCompile(
		`^[a-zA-Z0-9]{1}[a-zA-Z0-9\-]{1,31}$`,
	)

	// https://tools.ietf.org/html/rfc2141
	// <reserved> chars not included
	nssRegexp = regexp.MustCompile(
		`^(?:%[0-9A-Fa-f]{2}|[a-zA-Z0-9\-+(),.:=@;$_!*'])+$`,
	)
)

// https://tools.ietf.org/html/rfc3406#section-4.3
const (
	minNIDLength = 3
	maxNIDLength = 32
	minNSSLength = 1
)

type URN struct {
	nid []byte
	nss []byte
}

func New(
	nid string,
	nss string,
) (*URN, error) {
	bNID := bytes.TrimSpace([]byte(nid))

	err := validateNID(bNID)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create URN, reason: %s",
			err.Error(),
		)
	}

	bNSS := bytes.TrimSpace([]byte(nss))

	err = validateNSS(bNSS)
	if err != nil {
		return nil, fmt.Errorf(
			"can't create URN, reason: %s",
			err.Error(),
		)
	}

	return &URN{
		nid: bNID,
		nss: bNSS,
	}, nil
}

func Parse(rawURN string) (*URN, error) {
	return parseRawURN(
		bytes.TrimSpace([]byte(rawURN)),
	)
}

func MustParse(rawURN string) *URN {
	u, err := parseRawURN(
		bytes.TrimSpace([]byte(rawURN)),
	)

	if err != nil {
		panic(err)
	}

	return u
}

func parseRawURN(rawURN []byte) (*URN, error) {
	validURNPartsCount := 3
	tokens := bytes.SplitN(rawURN, urnDelimiter, validURNPartsCount)

	if len(tokens) != validURNPartsCount {
		return nil, fmt.Errorf(
			"invalid URN format, should be urn:<nid>:<nss>",
		)
	}

	prefix := tokens[0]
	if !bytes.Equal(prefix, urnPrefix) {
		return nil, fmt.Errorf(
			"URN '%s' must have prefix - %s",
			rawURN,
			urnPrefix,
		)
	}

	nid := tokens[1]
	if !nidRegexp.Match(nid) {
		return nil, fmt.Errorf(
			"NID %s doesn't satisfy pattern: %s",
			nid,
			nidRegexp.String(),
		)
	}

	nss := tokens[2]
	if !nssRegexp.Match(nss) {
		return nil, fmt.Errorf(
			"NSS %s doesn't satisfy the regexp rule: %s",
			nss,
			nssRegexp,
		)
	}

	return &URN{nid: nid, nss: nss}, nil
}

// String - returns string representation of a URN
func (urn *URN) String() string {
	return string(
		urn.constructURN(),
	)
}

// MarshalJSON - implements JSON Marshaller interface.
// Returns a valid JSON string
func (urn *URN) MarshalJSON() ([]byte, error) {
	return append(
		[]byte{34}, append(urn.constructURN(), []byte{34}...)...,
	), nil
}

// constructURN - constructs an URN in valid representation
// e.g. 'urn:<nid>:<nss>'
func (urn *URN) constructURN() []byte {
	return bytes.Join(
		[][]byte{
			urnPrefix,
			urn.nid,
			urn.nss,
		},
		urnDelimiter,
	)
}

func validateNID(nid []byte) error {
	if len(nid) < minNIDLength {
		return errors.New(
			"length of NID must be more than 2 letters long",
		)
	}

	if len(nid) > maxNIDLength {
		return fmt.Errorf(
			"NID must be mot greater than %d letters long",
			maxNIDLength,
		)
	}

	if bytes.Equal(bytes.ToLower(nid[:2]), experimentalNIDPrefix) {
		return fmt.Errorf(
			"NID %s is experimental",
			experimentalNIDPrefix,
		)
	}

	if bytes.Equal(bytes.ToLower(nid[:3]), xyNIDPrefix) {
		return fmt.Errorf(
			"NID %s mustn't start with: %s",
			nid,
			xyNIDPrefix,
		)
	}

	if bytes.Equal(bytes.ToLower(nid[:4]), reservedNIDPrefix) {
		return fmt.Errorf(
			"NID %s is reserved",
			reservedNIDPrefix,
		)
	}

	if !nidRegexp.Match(nid) {
		return fmt.Errorf(
			"NID %s doesn't satisfy pattern: %s",
			nid,
			nidRegexp.String(),
		)
	}

	return nil
}

func validateNSS(nss []byte) error {
	if len(nss) < minNSSLength {
		return fmt.Errorf(
			"NSS must be at least one character long",
		)
	}

	if !nssRegexp.Match(nss) {
		return fmt.Errorf(
			"NSS doesn't satisfy the regexp rule: %s",
			nssRegexp,
		)
	}

	return nil
}
