package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/michaelhoman/ShotSeek/internal/config"
	"github.com/michaelhoman/ShotSeek/internal/env"
)

type JWTService struct {
	config config.Config
	secret string
	expiry time.Duration
}

type JWTAuth struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func NewJWTService(secret string, expiry time.Duration) *JWTService {
	return &JWTService{secret: secret, expiry: expiry}
}

// NewJWTAuth initializes the JWTAuth struct by reading the ECDSA keys.
func NewJWTAuth() (*JWTAuth, error) {
	privateKey, err := loadPrivateKey(env.GetString("JWT_ECDSA_PRIVATE_KEY_PATH", ".keys/private_key.pem"))
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	publicKey, err := loadPublicKey(env.GetString("JWT_ECDSA_PUBLIC_KEY_PATH", ".keys/public_key.pem"))
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %v", err)
	}

	return &JWTAuth{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// loadPrivateKey reads and parses the private key from the file system.
func loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open private key file: %v", err)
	}
	defer file.Close()

	block, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	pemBlock, _ := pem.Decode(block)
	if pemBlock == nil {
		return nil, errors.New("failed to decode PEM block containing the private key")
	}

	privKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privKey, nil
}

// loadPublicKey reads and parses the public key from the file system.
func loadPublicKey(path string) (*ecdsa.PublicKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open public key file: %v", err)
	}
	defer file.Close()

	block, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %v", err)
	}

	pemBlock, _ := pem.Decode(block)
	if pemBlock == nil {
		return nil, errors.New("failed to decode PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an ECDSA public key")
	}

	return ecdsaPubKey, nil
}

// Function to create JWT
// func (a *AuthHandler) GenerateJWT(userID string, fingerprint string) (string, error) {
// 	// Get the secret key from the environment variable
// 	jwtSigningKey := []byte(os.Getenv("JWT_SIGNING_KEY"))
// 	// Set the expiration time (e.g., 1 hour from now)
// 	expirationTime := time.Now().Add(a.Config.Auth.Token.Exp).Unix()

// 	// Create the claims
// 	claims := Claims{
// 		Fingerprint: fingerprint, // Custom claim
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Issuer:    a.Config.Auth.Token.Iss,
// 			Audience:  jwt.ClaimStrings{a.Config.Auth.Token.Aud},
// 			Subject:   userID,
// 			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
// 			ExpiresAt: &jwt.NumericDate{Time: time.Unix(expirationTime, 0)},
// 		},
// 	}

// 	// Create the token using ES256 algorithm
// 	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

// 	// Sign the token with your secret key
// 	signedToken, err := token.SignedString(jwtSigningKey)
// 	if err != nil {
// 		return "", fmt.Errorf("could not sign the token: %v", err)
// 	}
// 	return signedToken, nil
// }

func (a *AuthHandler) GenerateJWTWithFP(userID uuid.UUID, fingerprint string) (string, error) {
	fmt.Println("321Generating JWT for user:", userID) // TODO REMOVE Debugging
	// Load the private key for signing (this can be done using the method we defined earlier)

	// privateKey, err := loadPrivateKey(a.JWTAuth.PrivateKey) // Path to your private key
	// if err != nil {
	// 	return "", fmt.Errorf("failed to load private key: %v", err)
	// }

	// Set the expiration time (e.g., 1 hour from now)
	expirationTime := time.Now().Add(a.Config.Auth.Token.Exp).Unix()

	// Create the claims
	claims := Claims{
		Fingerprint: fingerprint, // Custom claim
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.Config.Auth.Token.Iss,
			Audience:  jwt.ClaimStrings{a.Config.Auth.Token.Aud},
			Subject:   userID.String(),
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Unix(expirationTime, 0)},
		},
	}

	// Create the token using ES256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	fmt.Println("Token created:", token) // TODO REMOVE Debugging

	// Sign the token with your ECDSA private key
	signedToken, err := token.SignedString(a.JWTAuth.PrivateKey)
	fmt.Println("Signed token:", signedToken) // TODO REMOVE Debugging
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %v", err)
	}
	return signedToken, nil
}
func (a *AuthHandler) GenerateJWT(userID uuid.UUID) (string, error) {

	// Load the private key for signing (this can be done using the method we defined earlier)

	// privateKey, err := loadPrivateKey(a.JWTAuth.PrivateKey) // Path to your private key
	// if err != nil {
	// 	return "", fmt.Errorf("failed to load private key: %v", err)
	// }

	// Set the expiration time (e.g., 1 hour from now)
	expirationTime := time.Now().Add(a.Config.Auth.Token.Exp).Unix()

	// Create the claims
	claims := Claims{
		Fingerprint: "", // Custom claim
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.Config.Auth.Token.Iss,
			Audience:  jwt.ClaimStrings{a.Config.Auth.Token.Aud},
			Subject:   userID.String(),
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ExpiresAt: &jwt.NumericDate{Time: time.Unix(expirationTime, 0)},
		},
	}

	// Create the token using ES256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	fmt.Println("Token created:", token) // TODO REMOVE Debugging

	// Sign the token with your ECDSA private key
	signedToken, err := token.SignedString(a.JWTAuth.PrivateKey)
	fmt.Println("Signed token:", signedToken) // TODO REMOVE Debugging
	if err != nil {
		return "", fmt.Errorf("could not sign the token: %v", err)
	}
	return signedToken, nil
}

// Function to mask an IP address (only keep first two octets)
func anonymizeIP(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "" // Return empty string or handle error if needed
	}

	// Handle IPv4 addresses (e.g., "192.168.1.123" -> "192.168.0.0")
	if parsedIP.To4() != nil {
		// Split the IP into parts and mask the last two octets
		parts := strings.Split(ip, ".")
		if len(parts) >= 2 {
			return parts[0] + "." + parts[1] + ".0.0"
		}
	}

	// Handle IPv6 or return original IP for unsupported formats
	// Optionally, anonymize IPv6 addresses similarly
	if parsedIP.To16() != nil {
		// For simplicity, we can just return the first 4 parts of an IPv6 address
		// Masking the latter parts (IPv6 anonymization logic can vary based on requirements)
		parts := strings.Split(ip, ":")
		if len(parts) >= 4 {
			return strings.Join(parts[:4], ":") + "::"
		}
	}

	// Return the IP as-is if not an IPv4 or IPv6
	return ip
}

// Function to generate an IP fingerprint hash
func (a *AuthHandler) GenerateFingerprint(ip, userAgent string) string {
	// Anonymize IP before generating fingerprint
	anonymizedIP := anonymizeIP(ip)

	// Handle cases where user agent might be empty
	if userAgent == "" {
		userAgent = "unknown" // Default if empty
	}

	// Generate hash from anonymized IP + user agent
	data := anonymizedIP + userAgent
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]) // Convert hash to hex string
}

// Middleware to validate JWT fingerprint
func (a *AuthHandler) ValidateFingerprint(r *http.Request, expectedHash string) bool {
	// Get client IP (use X-Forwarded-For if available, fallback to RemoteAddr)
	// ip := r.Header.Get("X-Forwarded-For")
	ip := a.GetIPAddress(r)
	if ip == "" {
		ip = r.RemoteAddr
	}

	// // Ensure we only take the first IP in the X-Forwarded-For header (if multiple)
	// ips := strings.Split(ip, ",")
	// ip = strings.TrimSpace(ips[0])

	// Get the user agent from the request
	userAgent := r.UserAgent()

	// Generate current fingerprint and compare with expected fingerprint
	currentHash := a.GenerateFingerprint(ip, userAgent)
	return currentHash == expectedHash
}

func (a *AuthHandler) ValidateJWT(r *http.Request, tokenString, requestFingerprint string) (*Claims, error) {
	fmt.Println("Validating JWT Called: tokenString:", tokenString, "\nrequestFingerprint: ", requestFingerprint) // TODO REMOVE Debugging
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate Algorithm
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// // Get the secret key from the environment variable
		// jwtSigningKey := []byte(os.Getenv("JWT_SIGNING_KEY"))
		// return jwtSigningKey, nil
		return a.JWTAuth.PublicKey, nil

	},
		jwt.WithExpirationRequired(),                                // Ensure expiration is required and checked
		jwt.WithAudience(a.Config.Auth.Token.Aud),                   // Validate audience
		jwt.WithIssuer(a.Config.Auth.Token.Iss),                     // Validate issuer
		jwt.WithValidMethods([]string{jwt.SigningMethodES256.Name}), // Validate signing method
	)

	if err != nil {
		return nil, fmt.Errorf("token parsing error: %v", err)
	}

	// Validate Claims (Issuer, Fingerprint, etc.)
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Validate Issuer
	if claims.RegisteredClaims.Issuer != a.Config.Auth.Token.Iss {
		return nil, fmt.Errorf("invalid token issuer")
	}

	fmt.Printf("Extracted JWT Fingerprint: %s\n", claims.Fingerprint)
	fmt.Printf("Request Fingerprint: %s\n", requestFingerprint)
	// Validate Fingerprint
	if !a.ValidateFingerprint(r, claims.Fingerprint) {
		return nil, fmt.Errorf("invalid fingerprint")
	}

	// Validate Expiration
	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, fmt.Errorf("token is expired")
	}

	// Return the validated claims
	return claims, nil
}

// Function to extract JWT from the cookie
func GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", fmt.Errorf("could not get token from cookie: %v", err)
	}
	return cookie.Value, nil
}

// Helper function to extract the IP address
func (a *AuthHandler) GetIPAddress(r *http.Request) string {
	// Check if the IP is set by a proxy (e.g., Nginx, Cloudflare)
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		// Fallback to RemoteAddr if no proxy is used
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func (a *AuthHandler) authenticateRequest(r *http.Request) (*Claims, error) {
	// Step 1: Extract the JWT token from the cookie
	tokenString, err := GetTokenFromCookie(r)
	if err != nil {
		return nil, fmt.Errorf("token extraction failed: %v", err) // This could be a 401 Unauthorized in a real API
	}

	// Step 2: Validate the JWT token
	claims, err := a.ValidateJWT(r, tokenString, "")
	if err != nil {
		return nil, fmt.Errorf("JWT validation failed: %v", err)
	}

	// If token is valid, return claims
	return claims, nil
}

// Generate a secure random refresh token
func (a *AuthHandler) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 256-bit token
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil // Safe for URLs
}

// Hash the refresh token for storage
func (a *AuthHandler) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(hash[:]) // Store in DB
}

func (a *AuthHandler) ValidateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken, clientFP string) bool {
	storedTokens, err := a.store.Tokens.GetRefreshTokens(ctx, userID) // Fetch stored token and fingerprint                     // Hash the stored token for comparison
	if err != nil {
		return false
	}

	// Loop through all stored tokens and check if any match
	for _, token := range storedTokens {
		// Compare the token hash and fingerprint
		if token.TokenHash == refreshToken && token.StoredFP == clientFP {
			return true // Found a valid match
		}
	}

	// No matching token found
	return false
}

// Getter method for TokenExp
func (a *AuthHandler) GetAuthTokenExpiration() time.Duration {
	return a.Config.Auth.Token.Exp
}

// Getter method for TokenExp
func (a *AuthHandler) GetRefreshTokenExpiration() time.Duration {
	return a.Config.Auth.RefreshToken.Exp
}
