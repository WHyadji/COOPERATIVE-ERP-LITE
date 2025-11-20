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

// JWTAnggotaClaims mendefinisikan claims untuk JWT token anggota portal
type JWTAnggotaClaims struct {
	IDAnggota    uuid.UUID `json:"idAnggota"`
	IDKoperasi   uuid.UUID `json:"idKoperasi"`
	NomorAnggota string    `json:"nomorAnggota"`
	NamaLengkap  string    `json:"namaLengkap"`
	TipeToken    string    `json:"tipeToken"` // "anggota" untuk membedakan dari token pengguna
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

	// Validasi claims yang diperlukan
	if err := j.validateClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// validateClaims memvalidasi bahwa semua claims yang diperlukan ada dan valid
func (j *JWTUtil) validateClaims(claims *JWTClaims) error {
	// Validasi ID Pengguna tidak kosong
	if claims.IDPengguna == uuid.Nil {
		return errors.New("IDPengguna tidak boleh kosong")
	}

	// Validasi ID Koperasi tidak kosong
	if claims.IDKoperasi == uuid.Nil {
		return errors.New("IDKoperasi tidak boleh kosong")
	}

	// Validasi Peran tidak kosong
	if claims.Peran == "" {
		return errors.New("Peran tidak boleh kosong")
	}

	// Validasi Peran adalah salah satu dari yang diizinkan
	validRoles := map[models.PeranPengguna]bool{
		models.PeranAdmin:     true,
		models.PeranBendahara: true,
		models.PeranKasir:     true,
		models.PeranAnggota:   true,
	}

	if !validRoles[claims.Peran] {
		return errors.New("Peran tidak valid")
	}

	// Validasi NamaPengguna tidak kosong
	if claims.NamaPengguna == "" {
		return errors.New("NamaPengguna tidak boleh kosong")
	}

	return nil
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

// GenerateTokenAnggota menghasilkan JWT token untuk anggota portal
func (j *JWTUtil) GenerateTokenAnggota(anggota *models.Anggota) (string, error) {
	// Buat claims untuk anggota
	claims := JWTAnggotaClaims{
		IDAnggota:    anggota.ID,
		IDKoperasi:   anggota.IDKoperasi,
		NomorAnggota: anggota.NomorAnggota,
		NamaLengkap:  anggota.NamaLengkap,
		TipeToken:    "anggota",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "cooperative-erp-lite-portal",
			Subject:   anggota.ID.String(),
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

// ValidateTokenAnggota memvalidasi JWT token anggota dan mengembalikan claims
func (j *JWTUtil) ValidateTokenAnggota(tokenString string) (*JWTAnggotaClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTAnggotaClaims{}, func(token *jwt.Token) (interface{}, error) {
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
	claims, ok := token.Claims.(*JWTAnggotaClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid")
	}

	// Validasi tipe token
	if claims.TipeToken != "anggota" {
		return nil, errors.New("tipe token tidak valid")
	}

	// Validasi claims yang diperlukan
	if err := j.validateAnggotaClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// validateAnggotaClaims memvalidasi bahwa semua claims anggota yang diperlukan ada dan valid
func (j *JWTUtil) validateAnggotaClaims(claims *JWTAnggotaClaims) error {
	// Validasi ID Anggota tidak kosong
	if claims.IDAnggota == uuid.Nil {
		return errors.New("IDAnggota tidak boleh kosong")
	}

	// Validasi ID Koperasi tidak kosong
	if claims.IDKoperasi == uuid.Nil {
		return errors.New("IDKoperasi tidak boleh kosong")
	}

	// Validasi Nomor Anggota tidak kosong
	if claims.NomorAnggota == "" {
		return errors.New("NomorAnggota tidak boleh kosong")
	}

	// Validasi Tipe Token
	if claims.TipeToken != "anggota" {
		return errors.New("TipeToken harus 'anggota'")
	}

	return nil
}
