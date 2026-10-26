package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorive/coursepack/internal/cache"
	"github.com/victorive/coursepack/internal/cache/redis"
	"github.com/victorive/coursepack/internal/config"
	"github.com/victorive/coursepack/internal/email"
	"github.com/victorive/coursepack/internal/email/smtp"
	"github.com/victorive/coursepack/internal/email/templates"
	"github.com/victorive/coursepack/internal/feature/auth"
	"github.com/victorive/coursepack/internal/feature/class"
	"github.com/victorive/coursepack/internal/feature/course"
	"github.com/victorive/coursepack/internal/feature/department"
	"github.com/victorive/coursepack/internal/feature/file"
	"github.com/victorive/coursepack/internal/feature/semester"
	"github.com/victorive/coursepack/internal/feature/university"
	"github.com/victorive/coursepack/internal/feature/user"
	"github.com/victorive/coursepack/internal/permissions"
	jwt "github.com/victorive/coursepack/internal/pkg/jwt"
	"github.com/victorive/coursepack/internal/storage"
	"github.com/victorive/coursepack/internal/storage/local"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Container struct {
	Config             *config.Config
	Logger             *slog.Logger
	DB                 *pgxpool.Pool
	GORM               *gorm.DB
	JWT                *jwt.Manager
	Cache              cache.Cache
	Storage            storage.Storage
	Mailer             email.Mailer
	PermissionsService permissions.Service
	AuthService        auth.Service
	AuthRepository     auth.Repository
	UserRepository     user.Repository
	UserService        user.Service
	UniversityService  university.Service
	DepartmentService  department.Service
	ClassService       class.Service
	AcademicService    semester.Service
	CourseService      course.Service
	FileService        file.Service
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Container, error) {
	c := &Container{Config: cfg, Logger: logger}

	if err := c.initDB(ctx); err != nil {
		return nil, err
	}
	if err := c.initCache(); err != nil {
		return nil, err
	}
	if err := c.initStorage(); err != nil {
		return nil, err
	}
	if err := c.initMailer(); err != nil {
		return nil, err
	}

	c.JWT = jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpirationMinutes)
	c.initServices()

	return c, nil
}

func (c *Container) initDB(ctx context.Context) error {
	poolCfg, err := pgxpool.ParseConfig(c.Config.GetDSN())
	if err != nil {
		return fmt.Errorf("parse pgx config: %w", err)
	}
	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("new pgx pool: %w", err)
	}
	c.DB = db

	gormDB, err := gorm.Open(postgres.Open(c.Config.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}
	c.GORM = gormDB
	return nil
}

func (c *Container) initCache() error {
	r := c.Config.Redis
	redisCache, err := redis.New(r.Host, r.Port, r.Password, r.DB, r.PoolSize)
	if err != nil {
		return fmt.Errorf("redis.New: %w", err)
	}
	c.Cache = redisCache
	return nil
}

func (c *Container) initStorage() error {
	fileStorage, err := local.NewLocalStorage(c.Config.Storage.LocalPath)
	if err != nil {
		return fmt.Errorf("local.NewLocalStorage: %w", err)
	}
	c.Storage = fileStorage
	return nil
}

func (c *Container) initMailer() error {
	s := c.Config.SMTP
	mailer, err := smtp.New(s.Host, s.Port, s.User, s.Password, s.FromEmail, s.FromName, c.Config.App.AppURL, templates.FS)
	if err != nil {
		return fmt.Errorf("smtp.New: %w", err)
	}
	c.Mailer = mailer
	return nil
}

func (c *Container) initServices() {
	cfg := c.Config

	googleRedirectURL := cfg.Google.RedirectURL
	if googleRedirectURL == "" {
		googleRedirectURL = cfg.App.AppURL + "/api/v1/auth/google/callback"
	}

	authRepo := auth.NewRepository(c.GORM)
	c.AuthRepository = authRepo
	permissionsRepo := permissions.NewRepository(c.GORM)
	c.PermissionsService = permissions.NewService(permissionsRepo, c.Logger)
	c.AuthService = auth.NewService(
		authRepo,
		c.Cache,
		c.Mailer,
		c.JWT,
		auth.Config{
			OTPExpirationMinutes:     cfg.OTP.ExpirationMinutes,
			OTPMaxAttempts:           cfg.OTP.MaxAttempts,
			RefreshExpirationMinutes: cfg.JWT.RefreshExpirationMinutes,
			GoogleClientID:           cfg.Google.ClientID,
			GoogleClientSecret:       cfg.Google.ClientSecret,
			GoogleRedirectURL:        googleRedirectURL,
			FrontendURL:              cfg.App.FrontendURL,
		},
		c.PermissionsService,
		c.Logger,
	)

	c.UserRepository = user.NewRepository(c.GORM)
	c.UserService = user.NewService(
		user.NewRepository(c.GORM),
		c.Logger,
	)

	courseRepo := course.NewRepository(c.GORM)
	deptRepo := department.NewRepository(c.GORM)
	c.UniversityService = university.NewService(university.NewRepository(c.GORM))
	c.DepartmentService = department.NewService(deptRepo, c.Logger)
	c.ClassService = class.NewService(class.NewRepository(c.GORM), authRepo, deptRepo, c.PermissionsService, c.Logger)
	c.AcademicService = semester.NewService(semester.NewRepository(c.GORM), class.NewRepository(c.GORM), c.Logger)
	c.CourseService = course.NewService(courseRepo, class.NewRepository(c.GORM), semester.NewRepository(c.GORM), c.PermissionsService, c.Logger)
	pointRepo := file.NewPointRepository(c.GORM)
	c.FileService = file.NewService(file.NewRepository(c.GORM), authRepo, courseRepo, courseRepo, class.NewRepository(c.GORM), c.PermissionsService, c.Storage, cfg.Storage.MaxFileSize, pointRepo, c.Logger)
}

func (c *Container) Close() error {
	var firstErr error
	set := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if c.DB != nil {
		c.DB.Close()
	}
	if c.GORM != nil {
		if sqlDB, err := c.GORM.DB(); err == nil {
			set(sqlDB.Close())
		}
	}
	if c.Cache != nil {
		set(c.Cache.Close())
	}

	if firstErr != nil {
		c.Logger.Error("error during container close", "error", firstErr)
	}
	return firstErr
}
