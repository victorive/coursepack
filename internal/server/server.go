package server

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/victorive/coursepack/internal/app"
	"github.com/victorive/coursepack/internal/config"
	"github.com/victorive/coursepack/internal/pkg/apperr"
	"github.com/victorive/coursepack/internal/pkg/response"
)

type Server struct {
	router    *fiber.App
	config    *config.Config
	container *app.Container
	logger    *slog.Logger
}

func NewServer(cfg *config.Config, container *app.Container, log *slog.Logger) *Server {
	fapp := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ServerHeader: cfg.App.Name,
		ErrorHandler: errorHandler(log),
	})

	s := &Server{router: fapp, config: cfg, container: container, logger: log}
	s.setupMiddleware()
	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(requestid.New())

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowMethods:     s.config.CORS.AllowedMethods,
		AllowHeaders:     s.config.CORS.AllowedHeaders,
		AllowCredentials: false,
		MaxAge:           0,
	}))

	s.router.Use(logger.New(logger.Config{
		Format: "${time} | ${locals:requestid} | ${status} | ${latency} | ${ip} | ${method} | ${path}\n",
	}))

	s.router.Use(recover.New())

	s.router.Get("/health", s.healthCheck)

	if s.config.App.Env != "production" {
		s.router.Get("/docs/*", swagger.New(swagger.Config{
			DeepLinking:  true,
			DocExpansion: "list",
		}))
	}
}

func (s *Server) healthCheck(c *fiber.Ctx) error {
	return response.JSON(c, fiber.StatusOK, true, "OK", fiber.Map{
		"status": "healthy",
		"app":    s.config.App.Name,
	})
}

func (s *Server) RegisterRoutes() {
	RegisterRoutes(s.router, s.config, s.container)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.App.Host, s.config.App.Port)
	s.logger.Info("starting server", "addr", addr)
	return s.router.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.router.Shutdown()
}

func (s *Server) GetApp() *fiber.App {
	return s.router
}

const msgServerError = "something went wrong, please try again"

func errorHandler(log *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if appErr, ok := err.(*apperr.AppError); ok {
			if appErr.HTTPStatus >= fiber.StatusInternalServerError {
				log.Error("server error", "request_id", c.GetRespHeader(fiber.HeaderXRequestID), "path", c.Path(), "method", c.Method(), "error", appErr.Err)
				return c.Status(appErr.HTTPStatus).JSON(&response.Response{Success: false, Message: msgServerError})
			}

			if fe, ok := appErr.Data.(response.FieldErrors); ok {
				return c.Status(appErr.HTTPStatus).JSON(&response.ValidationErrorResponse{Success: false, Message: appErr.Message, Errors: fe})
			}

			return c.Status(appErr.HTTPStatus).JSON(&response.Response{Success: false, Message: appErr.Message})
		}

		code := fiber.StatusInternalServerError
		message := msgServerError

		if fiberErr, ok := err.(*fiber.Error); ok {
			code = fiberErr.Code
			if code < fiber.StatusInternalServerError {
				message = fiberErr.Message
			}
		}

		if code >= fiber.StatusInternalServerError {
			log.Error("unhandled error", "request_id", c.GetRespHeader(fiber.HeaderXRequestID), "status", code, "path", c.Path(), "method", c.Method(), "error", err)
		}

		return c.Status(code).JSON(&response.Response{Success: false, Message: message})
	}
}
