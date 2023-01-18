// Copyright © 2020 Ispirata Srl
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/astarte-platform/astarte-go/astarteservices"
	jwt "github.com/cristalhq/jwt/v3"
)

var (
	// ErrKeyMustBePEMEncoded is returned when the key is not encoded in PEM format
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded private key")
	// ErrNotPrivateKey is returned when the private key is not valid
	ErrNotPrivateKey = errors.New("Key is not a valid private key")
	// ErrUnsupportedPrivateKey is returned when the chosen private key is not supported for JWT generation
	ErrUnsupportedPrivateKey = errors.New("Key is not supported for JWT generation")
)

type AstarteClaims struct {
	jwt.StandardClaims

	AppEngineAPI    []string `json:"a_aea,omitempty"`
	Channels        []string `json:"a_ch,omitempty"`
	Flow            []string `json:"a_f,omitempty"`
	Housekeeping    []string `json:"a_ha,omitempty"`
	RealmManagement []string `json:"a_rma,omitempty"`
	Pairing         []string `json:"a_pa,omitempty"`
}

func (u *AstarteClaims) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// GenerateAstarteJWTFromKeyFile generates an Astarte Token for a specific API out of a Private Key File.
// servicesAndClaims specifies which services with which claims the token will be authorized to access. Leaving
// a claim empty will imply `.*::.*`, aka access to the entirety of the service's API tree
func GenerateAstarteJWTFromKeyFile(privateKeyFile string, servicesAndClaims map[astarteservices.AstarteService][]string,
	ttlSeconds int64) (jwtString string, err error) {
	keyPEM, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return "", err
	}

	return GenerateAstarteJWTFromPEMKey(keyPEM, servicesAndClaims, ttlSeconds)
}

// ParsePrivateKeyFromPEM parses a PEM encoded private key
func ParsePrivateKeyFromPEM(key []byte) (interface{}, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}

	case "PRIVATE KEY":
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}

	case "EC PRIVATE KEY":
		if parsedKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
			return nil, err
		}

	default:
		return nil, ErrNotPrivateKey
	}

	switch parsedKey.(type) {
	case *rsa.PrivateKey, *ecdsa.PrivateKey:
		return parsedKey, nil
	default:
		return nil, ErrUnsupportedPrivateKey
	}
}

// GenerateAstarteJWTFromPEMKey generates an Astarte Token for a specific API out of a Private Key PEM bytearray.
// servicesAndClaims specifies which services with which claims the token will be authorized to access. Leaving
// a claim empty will imply `.*::.*`, aka access to the entirety of the service's API tree
func GenerateAstarteJWTFromPEMKey(privateKeyPEM []byte, servicesAndClaims map[astarteservices.AstarteService][]string,
	ttlSeconds int64) (jwtString string, err error) {
	key, err := ParsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}

	// Build the token claims
	claims := AstarteClaims{}
	// Handle issue and expiry
	now := time.Now()
	claims.IssuedAt = jwt.NewNumericDate(now)
	if ttlSeconds > 0 {
		exp := now.Add(time.Duration(ttlSeconds) * time.Second)
		claims.ExpiresAt = jwt.NewNumericDate(exp)
	}

	for svc, c := range servicesAndClaims {
		if len(c) == 0 {
			switch svc {
			case astarteservices.Channels:
				c = []string{"JOIN::.*", "WATCH::.*"}
			default:
				c = []string{".*::.*"}
			}
		}

		switch svc {
		case astarteservices.AppEngine:
			claims.AppEngineAPI = c
		case astarteservices.Channels:
			claims.Channels = c
		case astarteservices.Flow:
			claims.Flow = c
		case astarteservices.Housekeeping:
			claims.Housekeeping = c
		case astarteservices.Pairing:
			claims.Pairing = c
		case astarteservices.RealmManagement:
			claims.RealmManagement = c
		}
	}

	signer, err := getJWTSigner(key)
	if err != nil {
		return "", err
	}
	builder := jwt.NewBuilder(signer)

	token, err := builder.Build(&claims)
	if err != nil {
		return "", err
	}

	return token.String(), nil
}

// GetJWTAstarteClaims returns the set of Astarte claims for an Astarte Token.
func GetJWTAstarteClaims(rawToken string) (AstarteClaims, error) {
	token, err := jwt.ParseString(rawToken)
	if err != nil {
		return AstarteClaims{}, err
	}

	ret := AstarteClaims{}
	err = json.Unmarshal(token.RawClaims(), &ret)
	if err != nil {
		return AstarteClaims{}, err
	}

	return ret, nil
}

// IsJWTAstarteClaimValidForService verifies that an Astarte Token has access to a given Astarte service.
func IsJWTAstarteClaimValidForService(token string, service astarteservices.AstarteService) (bool, error) {
	claims, err := GetJWTAstarteClaims(token)
	if err != nil {
		return false, err
	}
	switch service {
	case astarteservices.AppEngine:
		return hasAuth(claims.AppEngineAPI), nil
	case astarteservices.RealmManagement:
		return hasAuth(claims.RealmManagement), nil
	case astarteservices.Housekeeping:
		return hasAuth(claims.Housekeeping), nil
	case astarteservices.Pairing:
		return hasAuth(claims.Pairing), nil
	case astarteservices.Channels:
		return hasAuth(claims.Channels), nil
	case astarteservices.Flow:
		return hasAuth(claims.Flow), nil
	default:
		return false, fmt.Errorf("unknown Astarte service %s", service.String())
	}
}

func hasAuth(auth []string) bool {
	return len(auth) > 0
}

func getJWTSigner(key interface{}) (jwt.Signer, error) {
	var signer jwt.Signer
	var err error
	switch k := key.(type) {
	case *rsa.PrivateKey:
		signer, err = jwt.NewSignerRS(jwt.RS256, k)

	case *ecdsa.PrivateKey:
		// Match the EC curve with the correct signing algorithm
		switch k.PublicKey.Curve.Params().Name {
		case "P-256":
			signer, err = jwt.NewSignerES(jwt.ES256, k)
		case "P-384":
			signer, err = jwt.NewSignerES(jwt.ES384, k)
		case "P-521":
			signer, err = jwt.NewSignerES(jwt.ES512, k)
		default:
			return nil, ErrUnsupportedPrivateKey
		}
	}

	if err != nil {
		return nil, err
	}

	return signer, nil
}
