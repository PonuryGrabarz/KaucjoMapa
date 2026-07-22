package postgres

import (
	"database/sql"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

func toNullUUID(id *uuid.UUID) uuid.NullUUID {
	ID := uuid.NullUUID{
		UUID:  *id,
		Valid: id != nil,
	}

	return ID
}

func toNullString(str string) sql.NullString {
	nullStr := sql.NullString{
		String: str,
		Valid:  str != "",
	}

	return nullStr
}

func toNullFloat(num float64) sql.NullFloat64 {
	f := sql.NullFloat64{
		Float64: num,
		Valid:   num != 0,
	}

	return f
}

func toNullBool(b bool) sql.NullBool {
	nullBool := sql.NullBool{
		Bool:  b,
		Valid: true,
	}

	return nullBool
}

func toNullPointStatus(str entity.PointStatus) database.NullPointStatus {
	status := database.NullPointStatus{
		PointStatus: database.PointStatus(str),
		Valid:       str != "",
	}

	return status
}

func toNullPointType(str entity.PointType) database.NullPointType {
	pType := database.NullPointType{
		PointType: database.PointType(str),
		Valid:     str != "",
	}

	return pType
}

func toNullTime(t time.Time) sql.NullTime {
	nullTime := sql.NullTime{
		Time:  t,
		Valid: t != time.Time{},
	}

	return nullTime
}

func toCreateReturnPointParams(p entity.ReturnPoint) database.CreatePointParams {
	return database.CreatePointParams{
		ID:            p.ID,
		UserID:        toNullUUID(p.UserID),
		Name:          p.Name,
		Type:          database.PointType(p.Type),
		Longitude:     p.Longitude,
		Latitude:      p.Latitude,
		Address:       toNullString(p.Address),
		Pet:           p.Pet,
		Can:           p.Can,
		Glass:         p.Glass,
		OpeningHours:  toNullString(p.OpeningHours),
		Description:   toNullString(p.Description),
		CurrentStatus: database.PointStatus(p.CurrentStatus),
		IsVerified:    p.IsVerified,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func toReturnPointEntity(returnPoint database.ReturnPoint) entity.ReturnPoint {
	return entity.ReturnPoint{
		ID:            returnPoint.ID,
		UserID:        &returnPoint.UserID.UUID,
		Name:          returnPoint.Name,
		Type:          entity.PointType(returnPoint.Type),
		Longitude:     returnPoint.Longitude,
		Latitude:      returnPoint.Latitude,
		Address:       returnPoint.Address.String,
		Pet:           returnPoint.Pet,
		Can:           returnPoint.Can,
		Glass:         returnPoint.Glass,
		OpeningHours:  returnPoint.OpeningHours.String,
		Description:   returnPoint.Description.String,
		CurrentStatus: entity.PointStatus(returnPoint.CurrentStatus),
		IsVerified:    returnPoint.IsVerified,
		CreatedAt:     returnPoint.CreatedAt,
		UpdatedAt:     returnPoint.UpdatedAt,
	}
}

func toListReturnPointsParams(filters usecase.Filter) database.ListReturnPointsParams {
	return database.ListReturnPointsParams{
		Query:       toNullString(filters.Query),
		PointType:   toNullPointType(filters.Type),
		PointStatus: toNullPointStatus(filters.Status),
		Offset:      filters.Offset,
		Limit:       filters.Limit,
	}
}

func toCountReturnPointsParams(filters usecase.Filter) database.CountReturnPointsParams {
	return database.CountReturnPointsParams{
		Query:       toNullString(filters.Query),
		PointType:   toNullPointType(filters.Type),
		PointStatus: toNullPointStatus(filters.Status),
	}
}

func toUpdatePointParams(returnPoint entity.ReturnPoint) database.UpdatePointByIdParams {
	return database.UpdatePointByIdParams{
		ID:           returnPoint.ID,
		Name:         toNullString(returnPoint.Name),
		Type:         toNullPointType(returnPoint.Type),
		Address:      toNullString(returnPoint.Address),
		Latitude:     toNullFloat(returnPoint.Latitude),
		Longitude:    toNullFloat(returnPoint.Longitude),
		OpeningHours: toNullString(returnPoint.OpeningHours),
		Description:  toNullString(returnPoint.Description),
		Pet:          toNullBool(returnPoint.Pet),
		Can:          toNullBool(returnPoint.Can),
		Glass:        toNullBool(returnPoint.Glass),
	}
}

func toCreateUserParams(user entity.User) database.CreateUserParams {
	return database.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Reputation:   user.Reputation,
		IsAdmin:      user.IsAdmin,
		IsBanned:     user.IsBanned,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func toUserEntity(user database.User) entity.User {
	return entity.User{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Token:        "",
		Reputation:   user.Reputation,
		IsAdmin:      user.IsAdmin,
		IsBanned:     user.IsBanned,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func toCreateRefreshTokenParams(token entity.RefreshToken) database.CreateRefreshTokenParams {
	return database.CreateRefreshTokenParams{
		ID:        token.ID,
		Token:     token.Token,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
	}
}

func toGetUserRefresTokenParams(id uuid.UUID, tokenString string) database.GetUserRefreshTokenParams {
	return database.GetUserRefreshTokenParams{
		UserID: id,
		Token:  tokenString,
	}
}

func toRefreshTokenEntity(token database.RefreshToken) entity.RefreshToken {
	return entity.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: &token.RevokedAt.Time,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}
