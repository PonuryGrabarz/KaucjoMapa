package http

import (
	"os"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/handler"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/middleware"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/config"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

type FiberServer struct {
	*fiber.App
	ReturnPointSvc usecase.IReturnPointService
	UserSvc        usecase.IUserService
	AuthSvc        usecase.IAuthService
	ReportSvc      usecase.ICommunityReportsService
	VoteSvc        usecase.ICommunityReportVoteService
	cfg            *config.Config
}

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func New(returnPointSvc usecase.IReturnPointService, userSvc usecase.IUserService, authSvc usecase.IAuthService, reportSvc usecase.ICommunityReportsService, voteSvc usecase.ICommunityReportVoteService, tokenSvc port.ITokenSigner, hasherSvc port.IPasswordHasher, cfg *config.Config) *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader:    "KaucjoMapa",
			AppName:         "KaucjoMapa",
			ErrorHandler:    middleware.ErrorHandler,
			StructValidator: &structValidator{validate: validator.New()},
		}),
		ReturnPointSvc: returnPointSvc,
		UserSvc:        userSvc,
		AuthSvc:        authSvc,
		ReportSvc:      reportSvc,
		VoteSvc:        voteSvc,
		cfg:            cfg,
	}

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:4322", "http://127.0.0.1:3000", "http://localhost:4321", "https://localhost:3000", "https://localhost:3000", "https://kaucjomapa.pl", "https://www.kaucjomapa.pl"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"Set-Cookie",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	server.Use(recover.New())
	server.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))
	server.Use(helmet.New())
	server.Use(requestid.New())
	server.Use(logger.New(logger.Config{
		Format: "${requestid} [${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	api := server.Group("/api")
	authHandler := handler.NewAuthHandler(server.UserSvc, hasherSvc, tokenSvc, authSvc)
	returnPointHandler := handler.NewReturnPointHandler(server.ReturnPointSvc, server.UserSvc)
	userHandler := handler.NewUserHandler(server.UserSvc, tokenSvc, server.AuthSvc)
	reportsHandler := handler.NewCommunityReportHandler(server.ReportSvc, server.UserSvc)
	voteHandler := handler.NewCommunityReportVoteHandler(server.VoteSvc)

	ReturnPointRouter(api, *returnPointHandler, cfg.JWTSecret, tokenSvc)
	UsersRouter(api, *userHandler)
	AuthRouter(api, *authHandler, cfg.JWTSecret, tokenSvc)
	CommunityReportRouter(api, *reportsHandler, cfg.JWTSecret, tokenSvc)
	CommunityReportVoteRouter(api, *voteHandler, cfg.JWTSecret, tokenSvc)
	DocsRouter(api)

	return server
}
