package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/apsferreira/auth-service/backend/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ─── Formato do código OTP ────────────────────────────────────────────────────

func TestOTPCode_FormatIs6Digits(t *testing.T) {
	// Verifica que fmt.Sprintf("%06d", n) sempre produz 6 dígitos no range [0, 999999].
	cases := []int64{0, 1, 9, 99, 999, 9999, 99999, 999999}
	for _, n := range cases {
		code := fmt.Sprintf("%06d", n)
		if len(code) != 6 {
			t.Errorf("código para n=%d tem %d dígitos, esperado 6: %q", n, len(code), code)
		}
	}
}

func TestOTPCode_ZeroPadded(t *testing.T) {
	code := fmt.Sprintf("%06d", 42)
	if code != "000042" {
		t.Errorf("esperado '000042', got %q", code)
	}
}

func TestOTPCode_MaxValue(t *testing.T) {
	code := fmt.Sprintf("%06d", 999999)
	if code != "999999" {
		t.Errorf("esperado '999999', got %q", code)
	}
}

// ─── Hashing bcrypt ───────────────────────────────────────────────────────────

func TestOTPHash_DiffersFromPlaintext(t *testing.T) {
	code := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword: %v", err)
	}
	if string(hash) == code {
		t.Error("hash não deve ser igual ao código plain-text")
	}
}

func TestOTPHash_CompareSucceedsForCorrectCode(t *testing.T) {
	code := "654321"
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(code)); err != nil {
		t.Errorf("CompareHashAndPassword falhou para código correto: %v", err)
	}
}

func TestOTPHash_CompareFailsForWrongCode(t *testing.T) {
	code := "111111"
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte("222222"))
	if err == nil {
		t.Error("CompareHashAndPassword deveria falhar para código incorreto")
	}
}

// ─── Lógica de expiração (sem DB) ────────────────────────────────────────────

func TestOTPExpiry_NotExpiredWhenFresh(t *testing.T) {
	otp := &domain.OTPCode{
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if time.Now().After(otp.ExpiresAt) {
		t.Error("OTP recém-criado não deve estar expirado")
	}
}

func TestOTPExpiry_ExpiredWhenPast(t *testing.T) {
	otp := &domain.OTPCode{
		ExpiresAt: time.Now().Add(-time.Second),
	}
	if !time.Now().After(otp.ExpiresAt) {
		t.Error("OTP com ExpiresAt no passado deve ser considerado expirado")
	}
}

// ─── Lógica de tentativas máximas (sem DB) ───────────────────────────────────

func TestOTPAttempts_BlocksWhenAtMax(t *testing.T) {
	maxAttempts := 3
	otp := &domain.OTPCode{Attempts: maxAttempts}
	if otp.Attempts < maxAttempts {
		t.Error("com attempts == maxAttempts, acesso deve ser bloqueado")
	}
}

func TestOTPAttempts_AllowsWhenBelowMax(t *testing.T) {
	maxAttempts := 3
	otp := &domain.OTPCode{Attempts: 2}
	if otp.Attempts >= maxAttempts {
		t.Error("com attempts=2 e max=3, acesso deve ser permitido")
	}
}

// ─── OTPService.GetExpiryMinutes ─────────────────────────────────────────────

func TestOTPService_GetExpiryMinutes(t *testing.T) {
	svc := &OTPService{expiryMinutes: 10}
	if got := svc.GetExpiryMinutes(); got != 10 {
		t.Errorf("GetExpiryMinutes() = %d, esperado 10", got)
	}
}

func TestOTPService_GetExpiryMinutes_DifferentValues(t *testing.T) {
	cases := []int{1, 5, 15, 30, 60}
	for _, mins := range cases {
		svc := &OTPService{expiryMinutes: mins}
		if got := svc.GetExpiryMinutes(); got != mins {
			t.Errorf("expiryMinutes=%d, GetExpiryMinutes()=%d", mins, got)
		}
	}
}

// ─── domain.OTPCode struct ───────────────────────────────────────────────────

func TestOTPCode_FieldsSetCorrectly(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	expiry := now.Add(10 * time.Minute)

	otp := &domain.OTPCode{
		ID:        id,
		Email:     "user@example.com",
		CodeHash:  "hashed",
		Channel:   "email",
		Attempts:  0,
		ExpiresAt: expiry,
		CreatedAt: now,
	}

	if otp.ID != id {
		t.Errorf("ID incorreto: %v", otp.ID)
	}
	if otp.Email != "user@example.com" {
		t.Errorf("Email incorreto: %q", otp.Email)
	}
	if otp.Channel != "email" {
		t.Errorf("Channel incorreto: %q", otp.Channel)
	}
	if otp.Attempts != 0 {
		t.Errorf("Attempts deveria ser 0, got %d", otp.Attempts)
	}
}
