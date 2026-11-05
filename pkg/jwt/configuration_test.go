package jwt_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/buildbarn/bb-storage/pkg/jwt"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/require"
)

func TestJSONWebKeySetSignatureValidation(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	for _, test := range []struct {
		algorithm jose.SignatureAlgorithm
		key       crypto.Signer
	}{
		{jose.RS256, rsaKey},
		{jose.ES256, ecdsaKey},
		{jose.EdDSA, ed25519Key},
	} {
		t.Run(string(test.algorithm), func(t *testing.T) {
			keyID := "configured-key"
			encoded, err := json.Marshal(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{Key: test.key.Public(), KeyID: keyID}}})
			require.NoError(t, err)
			var keySet jose.JSONWebKeySet
			require.NoError(t, json.Unmarshal(encoded, &keySet))
			validator, err := jwt.NewSignatureValidatorFromJSONWebKeySet(&keySet)
			require.NoError(t, err)
			signer, err := jose.NewSigner(jose.SigningKey{Algorithm: test.algorithm, Key: test.key}, nil)
			require.NoError(t, err)
			signed, err := signer.Sign([]byte(`{"sub":"build-worker"}`))
			require.NoError(t, err)
			compact, err := signed.CompactSerialize()
			require.NoError(t, err)
			parts := strings.Split(compact, ".")
			signature, err := base64.RawURLEncoding.DecodeString(parts[2])
			require.NoError(t, err)
			input := parts[0] + "." + parts[1]
			require.True(t, validator.ValidateSignature(string(test.algorithm), &keyID, input, signature))
			require.False(t, validator.ValidateSignature(string(test.algorithm), &keyID, input+"x", signature))
			unknownKeyID := "unknown-key"
			require.False(t, validator.ValidateSignature(string(test.algorithm), &unknownKeyID, input, signature))
		})
	}
}
