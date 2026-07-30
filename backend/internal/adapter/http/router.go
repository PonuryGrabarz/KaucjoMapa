package http

import (
	"log"
	"os"
	"path/filepath"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/handler"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/middleware"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"
)

func ReturnPointRouter(app fiber.Router, returnPointHandler handler.ReturnPointHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Get("/return_points", returnPointHandler.GetReturnPoints)
	app.Post("/return_points", middleware.RequireJWT(auth, JWTSecret), returnPointHandler.Create)
	app.Delete("/return_points/:id", middleware.RequireJWT(auth, JWTSecret), returnPointHandler.DeletePoint)
	app.Put("/return_points/:id", middleware.RequireJWT(auth, JWTSecret), returnPointHandler.UpdatePoint)
}

func UsersRouter(app fiber.Router, userHandler handler.UserHandler, auth port.ITokenSigner, JWTSecret string) {
	app.Post("/users", userHandler.Create)
	app.Put("/users/password", middleware.RequireJWT(auth, JWTSecret), userHandler.ChangePassword)
}

func AuthRouter(app fiber.Router, authHandler handler.AuthHandler, oAuthHandler handler.OAuthHandler, recoveryHandler handler.PasswordRecoveryHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Post("/login", authHandler.Login)
	app.Post("/auth/token", middleware.ExtractJWT(auth, JWTSecret), authHandler.RefreshToken)
	app.Post("/auth/revoke", middleware.RequireJWT(auth, JWTSecret), authHandler.RevokeToken)
	app.Get("/auth/google", oAuthHandler.GoogleLoginHandler)
	app.Get("/auth/google/callback", oAuthHandler.GoogleCallbackHandler)
	app.Post("/auth/password-recovery", recoveryHandler.StartRecovery)
	app.Post("/auth/password-recovery/:challenge_id/verify", recoveryHandler.VerifyCode)
	app.Post("/auth/password-recovery/:challenge_id/resend", recoveryHandler.ResendCode)
	app.Post("/auth/password-recovery/:challenge_id/reset", recoveryHandler.ResetPassword)
}

func CommunityReportRouter(app fiber.Router, reportsHandler handler.CommunityReportHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Post("/reports", middleware.RequireJWT(auth, JWTSecret), reportsHandler.Create)
	app.Get("/reports/:return_point_id", reportsHandler.List)
	app.Delete("/reports/:id", middleware.RequireJWT(auth, JWTSecret), reportsHandler.Delete)
}

func CommunityReportVoteRouter(app fiber.Router, voteHandler handler.CommunityReportVoteHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Post("/votes/:report_id", middleware.RequireJWT(auth, JWTSecret), voteHandler.Create)
}

func NotificationsRouter(app fiber.Router, notificationsHandler handler.NotificationsHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Post("/notifications", middleware.RequireJWT(auth, JWTSecret), notificationsHandler.Create)
	app.Put("/notifications/:id", middleware.RequireJWT(auth, JWTSecret), notificationsHandler.MarkAsSeen)
	app.Get("/notifications", middleware.RequireJWT(auth, JWTSecret), notificationsHandler.List)
	app.Delete("/notifications/:id", middleware.RequireJWT(auth, JWTSecret), notificationsHandler.Delete)
}

func BookmarksRouter(app fiber.Router, bookmarkHandler handler.BookmarkHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Post("/bookmarks", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.Create)
	app.Get("/bookmarks", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.List)
	app.Get("/bookmarks/return_point/:return_point_id", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.GetByReturnPointID)
	app.Delete("/bookmarks/return_point/:return_point_id", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.DeleteByReturnPointID)
	app.Get("/bookmarks/:id", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.GetByID)
	app.Put("/bookmarks/:id", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.Update)
	app.Delete("/bookmarks/:id", middleware.RequireJWT(auth, JWTSecret), bookmarkHandler.Delete)
}

func ReturnPlansRouter(app fiber.Router, returnPlanHandler handler.ReturnPlanHandler, JWTSecret string, auth port.ITokenSigner) {
	app.Post("/return_plans", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.Create)
	app.Get("/return_plans", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.List)
	app.Get("/return_plans/draft", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.GetDraft)
	app.Get("/return_plans/:id", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.GetByID)
	app.Put("/return_plans/:id/targets", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.UpdateTargets)
	app.Post("/return_plans/:id/save", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.Save)
	app.Post("/return_plans/:id/confirm", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.ConfirmAll)
	app.Delete("/return_plans/:id", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.Delete)
	app.Post("/return_plans/:id/points", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.AddPoint)
	app.Put("/return_plans/:id/points/:point_id", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.UpdatePointQuantities)
	app.Post("/return_plans/:id/points/:point_id/confirm", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.ConfirmPoint)
	app.Delete("/return_plans/:id/points/:point_id", middleware.RequireJWT(auth, JWTSecret), returnPlanHandler.DeletePoint)
}

func DocsRouter(app fiber.Router) {
	openAPISepcDir, err := filepath.Abs("docs")
	if err != nil {
		log.Fatalf("Nie udało się znaleźć katalogu z plikiem OpenAPI: %v", err)
	}

	openAPISpec, err := os.ReadFile(filepath.Join(openAPISepcDir, "swagger.json"))
	if err != nil {
		log.Fatalf("Nie udało się wczytać pliku OpenAPI: %v", err)
	}

	app.Get("/docs/*", scalar.New(
		scalar.Config{
			BasePath:          "/api",
			Path:              "/docs",
			Title:             "KaucjoMapa API Documentation",
			FileContentString: string(openAPISpec),
		},
	))
}
