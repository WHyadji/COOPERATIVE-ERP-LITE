package utils

import (
	"cooperative-erp-lite/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims mendefinisikan claims untuk JWT token
type JWTClaims struct {
	IDPengguna   uuid.UUID            `json:"idPengguna"`
	IDKoperasi   uuid.UUID            `json:"idKoperasi"`
	NamaPengguna string               `json:"namaPengguna"`
	NamaLengkap  string               `json:"namaLengkap"`
	Peran        models.PeranPengguna `json:"peran"`
	jwt.RegisteredClaims
}

// JWTUtil adalah utility untuk JWT operations
type JWTUtil struct {
	secretKey       string
	expirationHours int
}

// NewJWTUtil membuat instance baru JWTUtil
func NewJWTUtil(secretKey string, expirationHours int) *JWTUtil {
	return &JWTUtil{
		secretKey:       secretKey,
		expirationHours: expirationHours,
	}
}

// GenerateToken menghasilkan JWT token untuk pengguna
func (j *JWTUtil) GenerateToken(pengguna *models.Pengguna) (string, error) {
	// Buat claims
	claims := JWTClaims{
		IDPengguna:   pengguna.ID,
		IDKoperasi:   pengguna.IDKoperasi,
		NamaPengguna: pengguna.NamaPengguna,
		NamaLengkap:  pengguna.NamaLengkap,
		Peran:        pengguna.Peran,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "cooperative-erp-lite",
			Subject:   pengguna.ID.String(),
		},
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token dan mengembalikan claims
func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing token tidak valid")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Ekstrak claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	return claims, nil
}

// RefreshToken menghasilkan token baru dari token yang sudah ada
func (j *JWTUtil) RefreshToken(tokenString string) (string, error) {
	// Validasi token lama
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Buat claims baru dengan expiration time yang baru
	newClaims := JWTClaims{
		IDPengguna:   claims.IDPengguna,
		IDKoperasi:   claims.IDKoperasi,
		NamaPengguna: claims.NamaPengguna,
		NamaLengkap:  claims.NamaLengkap,
		Peran:        claims.Peran,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "cooperative-erp-lite",
			Subject:   claims.IDPengguna.String(),
		},
	}

	// Buat token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	// Sign token dengan secret key
	newTokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}
