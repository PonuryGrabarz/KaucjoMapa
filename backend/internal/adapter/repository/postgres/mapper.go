package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

func toNullUUID(id *uuid.UUID) uuid.NullUUID {
	if id == nil {
		return uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	}

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

func toNullableFloat(value *float64) sql.NullFloat64 {
	if value == nil {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{Float64: *value, Valid: true}
}

func toNullableBool(value *bool) sql.NullBool {
	if value == nil {
		return sql.NullBool{}
	}

	return sql.NullBool{Bool: *value, Valid: true}
}

func toDatabaseRefundMethods(methods []entity.RefundMethod) []database.RefundMethod {
	if methods == nil {
		return nil
	}

	result := make([]database.RefundMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, database.RefundMethod(method))
	}

	return result
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

// func toNullTime(t time.Time) sql.NullTime {
// 	nullTime := sql.NullTime{
// 		Time:  t,
// 		Valid: t != time.Time{},
// 	}

// 	return nullTime
// }

func toCreateReturnPointParams(p entity.ReturnPoint) database.CreatePointParams {
	rawOpeningHours, _ := json.Marshal(p.OpeningHours)

	refundMethods := []database.RefundMethod{}
	for _, method := range p.RefundMethod {
		refundMethods = append(refundMethods, database.RefundMethod(method))
	}

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
		RefundMethods: refundMethods,
		OpeningHours:  rawOpeningHours,
		Description:   toNullString(p.Description),
		CurrentStatus: database.PointStatus(p.CurrentStatus),
		IsVerified:    p.IsVerified,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func toReturnPointEntity(returnPoint database.ReturnPoint) entity.ReturnPoint {
	var openingHours entity.OpeningHours
	json.Unmarshal(returnPoint.OpeningHours, &openingHours)

	refundMethods := []entity.RefundMethod{}
	for _, method := range returnPoint.RefundMethods {
		refundMethods = append(refundMethods, entity.RefundMethod(method))
	}

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
		RefundMethod:  refundMethods,
		OpeningHours:  openingHours,
		Description:   returnPoint.Description.String,
		CurrentStatus: entity.PointStatus(returnPoint.CurrentStatus),
		IsVerified:    returnPoint.IsVerified,
		CreatedAt:     returnPoint.CreatedAt,
		UpdatedAt:     returnPoint.UpdatedAt,
	}
}

func toListReturnPointEntity(returnPoint database.ListReturnPointsRow, includeDistance bool) entity.ReturnPoint {
	var openingHours entity.OpeningHours
	json.Unmarshal(returnPoint.OpeningHours, &openingHours)

	refundMethods := []entity.RefundMethod{}
	for _, method := range returnPoint.RefundMethods {
		refundMethods = append(refundMethods, entity.RefundMethod(method))
	}

	var distanceKm *float64
	if includeDistance {
		distanceKm = &returnPoint.DistanceKm
	}

	return entity.ReturnPoint{
		ID:            returnPoint.ID,
		UserID:        &returnPoint.UserID.UUID,
		Name:          returnPoint.Name,
		Type:          entity.PointType(returnPoint.Type),
		Longitude:     returnPoint.Longitude,
		Latitude:      returnPoint.Latitude,
		DistanceKm:    distanceKm,
		Address:       returnPoint.Address.String,
		Pet:           returnPoint.Pet,
		Can:           returnPoint.Can,
		Glass:         returnPoint.Glass,
		RefundMethod:  refundMethods,
		OpeningHours:  openingHours,
		Description:   returnPoint.Description.String,
		CurrentStatus: entity.PointStatus(returnPoint.CurrentStatus),
		IsVerified:    returnPoint.IsVerified,
		CreatedAt:     returnPoint.CreatedAt,
		UpdatedAt:     returnPoint.UpdatedAt,
	}
}

func toListReturnPointsParams(filters usecase.Filter) database.ListReturnPointsParams {
	return database.ListReturnPointsParams{
		UserID:        toNullUUID(filters.UserID),
		Query:         toNullString(filters.Query),
		PointType:     toNullPointType(filters.Type),
		PointStatus:   toNullPointStatus(filters.Status),
		Latitude:      toNullableFloat(filters.Latitude),
		Longitude:     toNullableFloat(filters.Longitude),
		MinDistanceKm: toNullableFloat(filters.MinDistanceKm),
		MaxDistanceKm: toNullableFloat(filters.MaxDistanceKm),
		Pet:           toNullableBool(filters.Pet),
		Can:           toNullableBool(filters.Can),
		Glass:         toNullableBool(filters.Glass),
		RefundMethods: toDatabaseRefundMethods(filters.RefundMethods),
		Offset:        filters.Offset,
		Limit:         filters.Limit,
	}
}

func toCountReturnPointsParams(filters usecase.Filter) database.CountReturnPointsParams {
	return database.CountReturnPointsParams{
		Query:         toNullString(filters.Query),
		PointType:     toNullPointType(filters.Type),
		PointStatus:   toNullPointStatus(filters.Status),
		Latitude:      toNullableFloat(filters.Latitude),
		Longitude:     toNullableFloat(filters.Longitude),
		MinDistanceKm: toNullableFloat(filters.MinDistanceKm),
		MaxDistanceKm: toNullableFloat(filters.MaxDistanceKm),
		Pet:           toNullableBool(filters.Pet),
		Can:           toNullableBool(filters.Can),
		Glass:         toNullableBool(filters.Glass),
		RefundMethods: toDatabaseRefundMethods(filters.RefundMethods),
	}
}

func toUpdatePointParams(returnPoint entity.ReturnPoint) database.UpdatePointByIdParams {
	rawOpeningHours, _ := json.Marshal(returnPoint.OpeningHours)

	nullRaw := pqtype.NullRawMessage{
		RawMessage: rawOpeningHours,
		Valid:      rawOpeningHours != nil,
	}

	// nil pozostawia stare metody zwrotu (COALESCE w SQL).
	var refundMethods []database.RefundMethod
	for _, method := range returnPoint.RefundMethod {
		refundMethods = append(refundMethods, database.RefundMethod(method))
	}

	return database.UpdatePointByIdParams{
		ID:            returnPoint.ID,
		Name:          toNullString(returnPoint.Name),
		Type:          toNullPointType(returnPoint.Type),
		Address:       toNullString(returnPoint.Address),
		Latitude:      toNullFloat(returnPoint.Latitude),
		Longitude:     toNullFloat(returnPoint.Longitude),
		OpeningHours:  nullRaw,
		Description:   toNullString(returnPoint.Description),
		Pet:           toNullBool(returnPoint.Pet),
		Can:           toNullBool(returnPoint.Can),
		Glass:         toNullBool(returnPoint.Glass),
		RefundMethods: refundMethods,
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
		ID:                 user.ID,
		Email:              user.Email,
		Username:           user.Username,
		PasswordHash:       user.PasswordHash,
		Token:              "",
		Reputation:         user.Reputation,
		TotalConfirmations: user.TotalConfirmations,
		TotalRejections:    user.TotalRejections,
		IsAdmin:            user.IsAdmin,
		IsBanned:           user.IsBanned,
		GoogleID:           user.GoogleSubject.String,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}
}

func toCreatePasswordResetChallengeParams(c entity.PasswordResetChallenge) database.CreatePasswordResetChallengeParams {
	return database.CreatePasswordResetChallengeParams{
		ID:                c.ID,
		UserID:            toNullUUID(c.UserID),
		CodeHash:          c.CodeHash,
		ExpiresAt:         c.ExpiresAt,
		ResendAvailableAt: c.ResendAvailableAt,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

func toUpdatePasswordResetChallengeCodeParams(c entity.PasswordResetChallenge) database.UpdatePasswordResetChallengeCodeParams {
	return database.UpdatePasswordResetChallengeCodeParams{
		ID:                c.ID,
		CodeHash:          c.CodeHash,
		ExpiresAt:         c.ExpiresAt,
		ResendAvailableAt: c.ResendAvailableAt,
	}
}

func toSetPasswordResetTokenParams(id uuid.UUID, tokenHash string, expiresAt time.Time) database.SetPasswordResetTokenParams {
	return database.SetPasswordResetTokenParams{
		ID:                  id,
		ResetTokenHash:      sql.NullString{String: tokenHash, Valid: true},
		ResetTokenExpiresAt: sql.NullTime{Time: expiresAt, Valid: true},
	}
}

func toPasswordResetChallengeEntity(c database.PasswordResetChallenge) entity.PasswordResetChallenge {
	var userID *uuid.UUID
	if c.UserID.Valid {
		id := c.UserID.UUID
		userID = &id
	}

	var resetTokenExpiresAt *time.Time
	if c.ResetTokenExpiresAt.Valid {
		resetTokenExpiresAt = &c.ResetTokenExpiresAt.Time
	}

	var consumedAt *time.Time
	if c.ConsumedAt.Valid {
		consumedAt = &c.ConsumedAt.Time
	}

	return entity.PasswordResetChallenge{
		ID:                  c.ID,
		UserID:              userID,
		CodeHash:            c.CodeHash,
		Attempts:            int(c.Attempts),
		ExpiresAt:           c.ExpiresAt,
		ResendAvailableAt:   c.ResendAvailableAt,
		ResetTokenHash:      c.ResetTokenHash.String,
		ResetTokenExpiresAt: resetTokenExpiresAt,
		ConsumedAt:          consumedAt,
		CreatedAt:           c.CreatedAt,
		UpdatedAt:           c.UpdatedAt,
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
	var revokedAt *time.Time
	if token.RevokedAt.Valid {
		revokedAt = &token.RevokedAt.Time
	}

	return entity.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		RevokedAt: revokedAt,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}

func toCreateCommunityReportParams(report entity.CommunityReport) database.CreateCommunityReportParams {
	return database.CreateCommunityReportParams{
		ID:            report.ID,
		UserID:        *report.UserID,
		ReturnPointID: *report.ReturnPointID,
		Status:        database.PointStatus(report.Status),
		Description:   report.Description,
	}
}

func toCommunityReportEntity(report database.CommunityReport) entity.CommunityReport {
	return entity.CommunityReport{
		ID:                 report.ID,
		UserID:             &report.UserID,
		ReturnPointID:      &report.ReturnPointID,
		Status:             entity.PointStatus(report.Status),
		Description:        report.Description,
		ConfirmationsCount: int(report.ConfirmationsCount),
		RejectionsCount:    int(report.RejectionsCount),
		CreatedAt:          report.CreatedAt,
		UpdatedAt:          report.UpdatedAt,
	}
}

func toListCommunityReportsByReturnPoint(row []database.ListCommunityReportsByReturnPointRow) []entity.UserCommunityReport {
	communityReports := []entity.UserCommunityReport{}

	for _, report := range row {
		communityReports = append(communityReports, entity.UserCommunityReport{
			Report: entity.CommunityReport{
				ID:                 report.ID,
				UserID:             &report.UserID,
				ReturnPointID:      &report.ReturnPointID,
				Status:             entity.PointStatus(report.Status),
				Description:        report.Description,
				ConfirmationsCount: int(report.ConfirmationsCount),
				RejectionsCount:    int(report.RejectionsCount),
				CreatedAt:          report.CreatedAt,
				UpdatedAt:          report.UpdatedAt,
			},
			Username:   report.Username,
			Reputation: report.Reputation,
		})
	}

	return communityReports
}

func toCreateCommunityReportVoteParams(vote entity.CommunityReportVote) database.UpsertCommunityReportVoteParams {
	return database.UpsertCommunityReportVoteParams{
		UserID:            *vote.UserID,
		CommunityReportID: *vote.CommunityReportID,
		Vote:              database.CommunityReportVoteValue(vote.Vote),
	}
}

func toCommunityReportVoteEntity(vote database.CommunityReportVote) entity.CommunityReportVote {
	return entity.CommunityReportVote{
		UserID:            &vote.UserID,
		CommunityReportID: &vote.CommunityReportID,
		Vote:              entity.CommunityReportVoteValue(vote.Vote),
		CreatedAt:         vote.CreatedAt,
		UpdatedAt:         vote.UpdatedAt,
	}
}

func toCreateNotificationParams(notification entity.Notification) database.CreateUserNotificationParams {
	return database.CreateUserNotificationParams{
		ID:          notification.ID,
		UserID:      toNullUUID(notification.UserID),
		Name:        notification.Name,
		Description: notification.Description,
		Seen:        notification.Seen,
		CreatedAt:   notification.CreatedAt,
	}
}

func toUpdateNotificationParams(notification entity.Notification, id uuid.UUID) database.UpdateUserNotificationParams {
	return database.UpdateUserNotificationParams{
		ID:          notification.ID,
		UserID:      toNullUUID(notification.UserID),
		Name:        notification.Name,
		Description: notification.Description,
		Seen:        notification.Seen,
	}

}

func toNotificationEntity(notification database.UserNotification) entity.Notification {
	return entity.Notification{
		ID:          notification.ID,
		UserID:      &notification.UserID.UUID,
		Name:        notification.Name,
		Description: notification.Description,
		Seen:        notification.Seen,
		CreatedAt:   notification.CreatedAt,
	}
}

func toGoogleUserParams(user entity.User) database.UpsertGoogleUserByEmailParams {
	return database.UpsertGoogleUserByEmailParams{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		GoogleSubject: toNullString(user.GoogleID),
	}
}

func toCreateBookmarkParams(bookmark entity.Bookmark) database.CreateBookmarkParams {
	return database.CreateBookmarkParams{
		ID:            bookmark.ID,
		UserID:        bookmark.UserID,
		ReturnPointID: bookmark.ReturnPointID,
		CreatedAt:     bookmark.CreatedAt,
	}
}

func toBookmarkEntity(bookmark database.Bookmark) entity.Bookmark {
	return entity.Bookmark{
		ID:            bookmark.ID,
		UserID:        bookmark.UserID,
		ReturnPointID: bookmark.ReturnPointID,
		CreatedAt:     bookmark.CreatedAt,
	}
}

func toBookmarksWithReturnPoints(rows []database.ListBookmarksWithReturnPointsByUserIdRow) []entity.BookmarkWithReturnPoint {
	result := make([]entity.BookmarkWithReturnPoint, 0, len(rows))
	for _, row := range rows {
		var openingHours entity.OpeningHours
		_ = json.Unmarshal(row.OpeningHours, &openingHours)

		refundMethods := make([]entity.RefundMethod, 0, len(row.RefundMethods))
		for _, method := range row.RefundMethods {
			refundMethods = append(refundMethods, entity.RefundMethod(method))
		}

		result = append(result, entity.BookmarkWithReturnPoint{
			Bookmark: entity.Bookmark{
				ID:            row.ID,
				UserID:        row.UserID,
				ReturnPointID: row.ReturnPointID,
				CreatedAt:     row.CreatedAt,
			},
			ReturnPoint: entity.ReturnPoint{
				ID:            row.ReturnPointID,
				Name:          row.Name,
				Type:          entity.PointType(row.Type),
				Address:       row.Address.String,
				Latitude:      row.Latitude,
				Longitude:     row.Longitude,
				Pet:           row.Pet,
				Can:           row.Can,
				Glass:         row.Glass,
				OpeningHours:  openingHours,
				Description:   row.Description.String,
				CurrentStatus: entity.PointStatus(row.CurrentStatus),
				IsVerified:    row.IsVerified,
				RefundMethod:  refundMethods,
			},
		})
	}

	return result
}
