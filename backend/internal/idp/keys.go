// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package idp

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"sync"
)

var (
	keyOnce    sync.Once
	privateKey *rsa.PrivateKey
	keyID      string
)

// ensureKey generates an RSA key pair once (in-memory; restart = new key)
func ensureKey() (*rsa.PrivateKey, error) {
	var err error
	keyOnce.Do(func() {
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return
		}
		keyID = "kkops-idp-1"
	})
	return privateKey, err
}

// getPrivateKey returns the signing key
func getPrivateKey() (*rsa.PrivateKey, error) {
	return ensureKey()
}

// getKeyID returns the key id for JWKS
func getKeyID() string {
	ensureKey()
	return keyID
}

// JWKS structure for discovery
type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// getJWKS returns the public key set as JSON for GET /.well-known/jwks.json
func getJWKS() ([]byte, error) {
	key, err := getPrivateKey()
	if err != nil {
		return nil, err
	}
	n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.PublicKey.E)).Bytes())
	j := jwks{
		Keys: []jwk{{
			Kty: "RSA",
			Kid: getKeyID(),
			Use: "sig",
			Alg: "RS256",
			N:   n,
			E:   e,
		}},
	}
	return json.Marshal(j)
}
