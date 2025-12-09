package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

type JWKSClient struct {
	jwksURL    string
	httpClient *http.Client
	cache      *JWKS
	cacheMu    sync.RWMutex
	cacheTime  time.Time
	cacheTTL   time.Duration
}

func NewJWKSClient(domain string) *JWKSClient {
	return &JWKSClient{
		jwksURL: fmt.Sprintf("https://%s/.well-known/jwks.json", domain),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cacheTTL: 1 * time.Hour,
	}
}

func (c *JWKSClient) GetKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	jwks, err := c.getJWKS(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return parseRSAPublicKey(key)
		}
	}

	c.invalidateCache()
	jwks, err = c.getJWKS(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return parseRSAPublicKey(key)
		}
	}

	return nil, fmt.Errorf("key %s not found in JWKS", kid)
}

func (c *JWKSClient) getJWKS(ctx context.Context) (*JWKS, error) {
	c.cacheMu.RLock()
	if c.cache != nil && time.Since(c.cacheTime) < c.cacheTTL {
		defer c.cacheMu.RUnlock()
		return c.cache, nil
	}
	c.cacheMu.RUnlock()

	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if c.cache != nil && time.Since(c.cacheTime) < c.cacheTTL {
		return c.cache, nil
	}

	jwks, err := c.fetchJWKS(ctx)
	if err != nil {
		return nil, err
	}

	c.cache = jwks
	c.cacheTime = time.Now()
	return jwks, nil
}

func (c *JWKSClient) fetchJWKS(ctx context.Context) (*JWKS, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	return &jwks, nil
}

func (c *JWKSClient) invalidateCache() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.cache = nil
}

func parseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	if jwk.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type: %s", jwk.Kty)
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}
