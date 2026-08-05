package http

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/handler"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/middleware"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
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
	ReturnPointSvc      usecase.IReturnPointService
	UserSvc             usecase.IUserService
	AuthSvc             usecase.IAuthService
	ReportSvc           usecase.ICommunityReportsService
	VoteSvc             usecase.ICommunityReportVoteService
	NotificationsSvc    usecase.INotificationsService
	BookmarkSvc         usecase.IBookmarkService
	OAuthSvc            usecase.IGoogleOAuthService
	PasswordRecoverySvc usecase.IPasswordRecoveryService
	ReturnPlanSvc       usecase.IReturnPlanService
	DataErrorReportSvc  usecase.IDataErrorReportsService
	HasherSvc           usecase.IPasswordHasher
	cfg                 *config.Config
}

type structValidator struct {
	validate *validator.Validate
}

// newValidator zwraca validator, w którym nazwy pól w błędach pochodzą z tagu
// `json` (a nie z nazwy pola w Go), dzięki czemu klient widzi znane mu klucze.
func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// Validate implementuje interfejs walidatora Fibera. Błędy walidacji są
// tłumaczone na entity.AppError (KindBadRequest) z listą pól, żeby warstwa
// HTTP mogła zwrócić czytelne szczegóły.
func (v *structValidator) Validate(out any) error {
	err := v.validate.Struct(out)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		fields := make([]entity.FieldError, 0, len(validationErrs))
		for _, fe := range validationErrs {
			fields = append(fields, entity.FieldError{
				Field:   fe.Field(),
				Rule:    fe.Tag(),
				Param:   fe.Param(),
				Message: translateFieldError(fe),
			})
		}
		return entity.NewBadRequest("Validation failed", fields...)
	}

	// Nietypowy błąd walidatora (np. zła konfiguracja) chowamy jako 500.
	return entity.NewInternal(err)
}

// translateFieldError buduje czytelny komunikat dla pojedynczej reguły.
func translateFieldError(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "oneof":
		return field + " must be one of: " + fe.Param()
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	case "min":
		return field + " must be at least " + fe.Param()
	case "email":
		return field + " must be a valid email"
	case "latitude", "longitude":
		return field + " must be a valid " + fe.Tag()
	default:
		return field + " is invalid (" + fe.Tag() + ")"
	}
}

func New(returnPointSvc usecase.IReturnPointService, userSvc usecase.IUserService, authSvc usecase.IAuthService, reportSvc usecase.ICommunityReportsService, voteSvc usecase.ICommunityReportVoteService, notSvc usecase.INotificationsService, bookmarkSvc usecase.IBookmarkService, returnPlanSvc usecase.IReturnPlanService, oAuthSvc usecase.IGoogleOAuthService, passwordRecoverySvc usecase.IPasswordRecoveryService, dataErrorReportsSvc usecase.IDataErrorReportsService, tokenSvc port.ITokenSigner, hasherSvc port.IPasswordHasher, cfg *config.Config) *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader:    "KaucjoMapa",
			AppName:         "KaucjoMapa",
			BodyLimit:       2 * 1024 * 1024,
			ErrorHandler:    middleware.ErrorHandler,
			StructValidator: &structValidator{validate: newValidator()},
		}),
		ReturnPointSvc:      returnPointSvc,
		UserSvc:             userSvc,
		AuthSvc:             authSvc,
		ReportSvc:           reportSvc,
		VoteSvc:             voteSvc,
		NotificationsSvc:    notSvc,
		BookmarkSvc:         bookmarkSvc,
		ReturnPlanSvc:       returnPlanSvc,
		PasswordRecoverySvc: passwordRecoverySvc,
		OAuthSvc:            oAuthSvc,
		HasherSvc:           hasherSvc,
		DataErrorReportSvc:  dataErrorReportsSvc,
		cfg:                 cfg,
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
	notificationsHandler := handler.NewNotificationsHandler(server.NotificationsSvc)
	bookmarkHandler := handler.NewBookmarkHandler(server.BookmarkSvc)
	returnPlanHandler := handler.NewReturnPlanHandler(server.ReturnPlanSvc)
	oAuthHandler := handler.NewOAuthHandler(server.UserSvc, server.HasherSvc, server.OAuthSvc)
	recoveryHandler := handler.NewPasswordRecoveryHandler(server.PasswordRecoverySvc)
	dataErrorReportHandler := handler.NewDataErrorReportHandler(server.DataErrorReportSvc)

	ReturnPointRouter(api, *returnPointHandler, cfg.JWTSecret, tokenSvc)
	UsersRouter(api, *userHandler, tokenSvc, cfg.JWTSecret)
	AuthRouter(api, *authHandler, *oAuthHandler, *recoveryHandler, cfg.JWTSecret, tokenSvc)
	CommunityReportRouter(api, *reportsHandler, cfg.JWTSecret, tokenSvc)
	CommunityReportVoteRouter(api, *voteHandler, cfg.JWTSecret, tokenSvc)
	NotificationsRouter(api, *notificationsHandler, cfg.JWTSecret, tokenSvc)
	BookmarksRouter(api, *bookmarkHandler, cfg.JWTSecret, tokenSvc)
	ReturnPlansRouter(api, *returnPlanHandler, cfg.JWTSecret, tokenSvc)
	DataErrorReportRouter(api, *dataErrorReportHandler, cfg.JWTSecret, tokenSvc)
	DocsRouter(api)

	return server
}
