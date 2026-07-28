package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/google/uuid"
)

const (
	recoveryCodeTTL     = 10 * time.Minute
	recoveryResendDelay = 60 * time.Second
	recoveryMaxAttempts = 5
	resetTokenTTL       = 10 * time.Minute
)

type IPasswordRecoveryService interface {
	StartRecovery(ctx context.Context, in dto.StartRecoveryInput) (entity.PasswordResetChallenge, error)
	VerifyCode(ctx context.Context, in dto.VerifyCodeInput) (string, time.Time, error)
	ResendCode(ctx context.Context, in dto.ResendCodeInput) (entity.PasswordResetChallenge, error)
	ResetPassword(ctx context.Context, in dto.ResetPasswordInput) error
}

type passwordRecoveryService struct {
	recoveryRepo     IPasswordRecoveryRepository
	userRepo         IUserRepository
	refreshTokenRepo IRefreshTokenRepository
	hasherSvc        port.IPasswordHasher
	mailSvc          port.IMailSender
}

func NewPasswordRecoveryService(recoveryRepo IPasswordRecoveryRepository, userRepo IUserRepository, refreshTokenRepo IRefreshTokenRepository, hasherSvc port.IPasswordHasher, mailSvc port.IMailSender) IPasswordRecoveryService {
	return &passwordRecoveryService{
		recoveryRepo:     recoveryRepo,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		hasherSvc:        hasherSvc,
		mailSvc:          mailSvc,
	}
}

func (s *passwordRecoveryService) StartRecovery(ctx context.Context, in dto.StartRecoveryInput) (entity.PasswordResetChallenge, error) {
	var userID *uuid.UUID
	user, err := s.userRepo.GetUserByEmail(ctx, in.Email)

	switch {
	case err == nil:
		userID = &user.ID
		if err := s.recoveryRepo.ConsumeActiveByUserID(ctx, user.ID); err != nil {
			return entity.PasswordResetChallenge{}, err
		}
	case entity.IsKind(err, entity.KindNotFound):
	default:
		return entity.PasswordResetChallenge{}, err
	}

	code, err := generateRecoveryCode()
	if err != nil {
		return entity.PasswordResetChallenge{}, entity.NewInternal(err)
	}

	codeHash, err := s.hasherSvc.HashPassword(code)
	if err != nil {
		return entity.PasswordResetChallenge{}, entity.NewInternal(err)
	}

	now := time.Now()
	challenge, err := s.recoveryRepo.Create(ctx, entity.PasswordResetChallenge{
		ID:                uuid.New(),
		UserID:            userID,
		CodeHash:          codeHash,
		ExpiresAt:         now.Add(recoveryCodeTTL),
		ResendAvailableAt: now.Add(recoveryResendDelay),
		CreatedAt:         now,
		UpdatedAt:         now,
	})
	if err != nil {
		return entity.PasswordResetChallenge{}, err
	}

	if userID != nil {
		if err := s.mailSvc.SendPasswordRecoveryCode(ctx, in.Email, code); err != nil {
			return entity.PasswordResetChallenge{}, entity.NewInternal(err)
		}
	}

	return challenge, nil
}

func (s *passwordRecoveryService) VerifyCode(ctx context.Context, in dto.VerifyCodeInput) (string, time.Time, error) {
	challenge, err := s.getActiveChallenge(ctx, in.ChallengeID)
	if err != nil {
		return "", time.Time{}, err
	}

	if challenge.UserID == nil || challenge.ResetTokenHash != "" {
		return "", time.Time{}, s.failAttempt(ctx, challenge)
	}

	match, err := s.hasherSvc.CheckPasswordHash(in.Code, challenge.CodeHash)
	if err != nil {
		return "", time.Time{}, entity.NewInternal(err)
	}
	if !match {
		return "", time.Time{}, s.failAttempt(ctx, challenge)
	}

	resetToken, tokenHash, err := generateResetToken()
	if err != nil {
		return "", time.Time{}, entity.NewInternal(err)
	}

	expiresAt := time.Now().Add(resetTokenTTL)
	_, err = s.recoveryRepo.SetResetToken(ctx, challenge.ID, tokenHash, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return resetToken, expiresAt, nil
}

func (s *passwordRecoveryService) ResendCode(ctx context.Context, in dto.ResendCodeInput) (entity.PasswordResetChallenge, error) {
	challenge, err := s.recoveryRepo.GetByID(ctx, in.ChallengeID)
	if err != nil {
		return entity.PasswordResetChallenge{}, entity.NewNotFound("challenge not found")
	}

	if challenge.ConsumedAt != nil {
		return entity.PasswordResetChallenge{}, entity.NewNotFound("challenge not found")
	}

	if time.Now().Before(challenge.ResendAvailableAt) {
		return entity.PasswordResetChallenge{}, entity.NewConflict("resend not yet available")
	}

	code, err := generateRecoveryCode()
	if err != nil {
		return entity.PasswordResetChallenge{}, entity.NewInternal(err)
	}

	codeHash, err := s.hasherSvc.HashPassword(code)
	if err != nil {
		return entity.PasswordResetChallenge{}, entity.NewInternal(err)
	}

	now := time.Now()
	challenge.CodeHash = codeHash
	challenge.ExpiresAt = now.Add(recoveryCodeTTL)
	challenge.ResendAvailableAt = now.Add(recoveryResendDelay)

	updated, err := s.recoveryRepo.UpdateCode(ctx, challenge)
	if err != nil {
		return entity.PasswordResetChallenge{}, err
	}

	if challenge.UserID != nil {
		user, err := s.userRepo.GetUserByID(ctx, *challenge.UserID)
		if err != nil {
			return entity.PasswordResetChallenge{}, err
		}
		if err := s.mailSvc.SendPasswordRecoveryCode(ctx, user.Email, code); err != nil {
			return entity.PasswordResetChallenge{}, entity.NewInternal(err)
		}
	}

	return updated, nil
}

func (s *passwordRecoveryService) ResetPassword(ctx context.Context, in dto.ResetPasswordInput) error {
	challenge, err := s.recoveryRepo.GetByID(ctx, in.ChallengeID)
	if err != nil {
		return entity.NewNotFound("challenge not found")
	}

	if challenge.ConsumedAt != nil || challenge.UserID == nil || challenge.ResetTokenHash == "" {
		return entity.NewUnauthorized("invalid reset token")
	}

	if challenge.ResetTokenExpiresAt == nil || time.Now().After(*challenge.ResetTokenExpiresAt) {
		return entity.NewUnauthorized("reset token expired")
	}

	tokenHash := hashResetToken(in.ResetToken)
	if subtle.ConstantTimeCompare([]byte(tokenHash), []byte(challenge.ResetTokenHash)) != 1 {
		return entity.NewUnauthorized("invalid reset token")
	}

	passwordHash, err := s.hasherSvc.HashPassword(in.NewPassword)
	if err != nil {
		return entity.NewInternal(err)
	}

	if err := s.userRepo.UpdatePassword(ctx, *challenge.UserID, passwordHash); err != nil {
		return err
	}

	if err := s.recoveryRepo.Consume(ctx, challenge.ID); err != nil {
		return err
	}

	// Po zmianie hasła wszystkie dotychczasowe sesje tracą ważność.
	if err := s.refreshTokenRepo.RevokeAllUserTokens(ctx, *challenge.UserID); err != nil {
		return err
	}

	return nil
}

func (s *passwordRecoveryService) getActiveChallenge(ctx context.Context, id uuid.UUID) (entity.PasswordResetChallenge, error) {
	challenge, err := s.recoveryRepo.GetByID(ctx, id)
	if err != nil {
		return entity.PasswordResetChallenge{}, entity.NewNotFound("challenge not found")
	}

	if challenge.ConsumedAt != nil || time.Now().After(challenge.ExpiresAt) {
		return entity.PasswordResetChallenge{}, entity.NewNotFound("challenge not found")
	}

	if challenge.Attempts >= recoveryMaxAttempts {
		return entity.PasswordResetChallenge{}, entity.NewForbidden("too many attempts")
	}

	return challenge, nil
}

func (s *passwordRecoveryService) failAttempt(ctx context.Context, challenge entity.PasswordResetChallenge) error {
	updated, err := s.recoveryRepo.IncrementAttempts(ctx, challenge.ID)
	if err != nil {
		return err
	}

	if updated.Attempts >= recoveryMaxAttempts {
		if err := s.recoveryRepo.Consume(ctx, challenge.ID); err != nil {
			return err
		}
		return entity.NewForbidden("too many attempts")
	}

	return entity.NewUnauthorized("invalid code")
}

func generateRecoveryCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(100_000_000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%08d", n.Int64()), nil
}

func generateResetToken() (token string, tokenHash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}

	token = base64.RawURLEncoding.EncodeToString(buf)
	return token, hashResetToken(token), nil
}

func hashResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
