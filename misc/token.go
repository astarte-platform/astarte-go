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

package misc

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"

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

type astarteClaims struct {
	jwt.StandardClaims

	AppEngineAPI    []string `json:"a_aea,omitempty"`
	Channels        []string `json:"a_ch,omitempty"`
	Flow            []string `json:"a_f,omitempty"`
	Housekeeping    []string `json:"a_ha,omitempty"`
	RealmManagement []string `json:"a_rma,omitempty"`
	Pairing         []string `json:"a_pa,omitempty"`
}

func (u *astarteClaims) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// GenerateAstarteJWTFromKeyFile generates an Astarte Token for a specific API out of a Private Key File.
// servicesAndClaims specifies which services with which claims the token will be authorized to access. Leaving
// a claim empty will imply `.*::.*`, aka access to the entirety of the service's API tree
func GenerateAstarteJWTFromKeyFile(privateKeyFile string, servicesAndClaims map[AstarteService][]string,
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
func GenerateAstarteJWTFromPEMKey(privateKeyPEM []byte, servicesAndClaims map[AstarteService][]string,
	ttlSeconds int64) (jwtString string, err error) {
	key, err := ParsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}

	// Build the token claims
	claims := astarteClaims{}
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
			case Channels:
				c = []string{"JOIN::.*", "WATCH::.*"}
			default:
				c = []string{".*::.*"}
			}
		}

		switch svc {
		case AppEngine:
			claims.AppEngineAPI = c
		case Channels:
			claims.Channels = c
		case Flow:
			claims.Flow = c
		case Housekeeping:
			claims.Housekeeping = c
		case Pairing:
			claims.Pairing = c
		case RealmManagement:
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
