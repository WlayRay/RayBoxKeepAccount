package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"ray_box/infrastructure/config"
	"ray_box/infrastructure/zlog"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	POSTGRESQL_DB_MAIN = "ray_box"
)

var (
	postgresConn  = make(map[string]*gorm.DB)
	postgresMutex sync.RWMutex
	pgLogger      = zlog.GetLogger()
)

func GetPostgresConn(db string) *gorm.DB {
	postgresMutex.RLock()
	conn, ok := postgresConn[db]
	postgresMutex.RUnlock()

	if !ok {
		postgresMutex.Lock()
		defer postgresMutex.Unlock()

		dbMap := map[string]string{
			POSTGRESQL_DB_MAIN: "PG_RAY_BOX",
		}
		envPrefix := dbMap[db]

		userName := config.GetConfig(envPrefix + "_USERNAME")
		userPwd := config.GetConfig(envPrefix + "_PASSWORD")
		host := config.GetConfig(envPrefix + "_HOST")
		port := config.GetConfig(envPrefix + "_PORT")

		// 设置日志级别
		envLogLevel, _ := os.LookupEnv("LIB_LOG_LEVEL")
		var gormlevel gormLogger.LogLevel
		switch envLogLevel {
		case "debug":
			gormlevel = gormLogger.Info
		case "info":
			gormlevel = gormLogger.Info
		case "warning":
			gormlevel = gormLogger.Warn
		case "error", "dpanic", "panic", "fatal":
			gormlevel = gormLogger.Error
		default:
			gormlevel = gormLogger.Warn
		}

		// 构建DSN字符串
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			host, userName, userPwd, db, port)

		dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: &CustomPostgresLogger{
				Logger: pgLogger,
				Config: gormLogger.Config{
					LogLevel:      gormlevel,
					SlowThreshold: 500 * time.Millisecond,
				},
			},
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			pgLogger.Error("Failed to connect to PostgreSQL",
				zap.Error(err),
				zap.String("database", db),
			)
			return nil
		} else {
			// 设置连接池
			if sqlDB, err := dbConn.DB(); err == nil {
				sqlDB.SetMaxOpenConns(20)
				sqlDB.SetMaxIdleConns(5)
				sqlDB.SetConnMaxLifetime(30 * time.Minute)
			}

			pgLogger.Debug("postgresql connection dsn " + dsn)

			if strings.ToLower(envLogLevel) == "debug" {
				postgresConn[db] = dbConn.Debug()
			} else {
				postgresConn[db] = dbConn
			}
		}
		conn = dbConn
	}

	return conn
}

type CustomPostgresLogger struct {
	Logger *zap.Logger
	Config gormLogger.Config
}

// LogMode log mode
func (l *CustomPostgresLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newlogger := *l
	newlogger.Config.LogLevel = level
	return &newlogger
}

// Info print info
func (l CustomPostgresLogger) Info(ctx context.Context, msg string, data ...any) {
	defer l.Logger.Sync()

	l.Logger.Error(fmt.Sprintf("%s", append([]any{msg}, data...)...),
		zap.String("source", utils.FileWithLineNum()),
		zap.String("agg_type", "gorm"),
	)
}

// Warn print warn messages
func (l CustomPostgresLogger) Warn(ctx context.Context, msg string, data ...any) {
	defer l.Logger.Sync()

	l.Logger.Error(fmt.Sprintf("%s", append([]any{msg}, data...)...),
		zap.String("source", utils.FileWithLineNum()),
		zap.String("agg_type", "gorm"),
	)
}

// Error print error messages
func (l CustomPostgresLogger) Error(ctx context.Context, msg string, data ...any) {
	defer l.Logger.Sync()

	l.Logger.Error(fmt.Sprintf("%s", append([]any{msg}, data...)...),
		zap.String("source", utils.FileWithLineNum()),
		zap.String("agg_type", "gorm"),
	)
}

// Trace print sql message
func (l CustomPostgresLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	defer l.Logger.Sync()

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Config.LogLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.Config.IgnoreRecordNotFoundError):
		sql, rows := fc()
		l.Logger.Error(err.Error(),
			zap.String("source", utils.FileWithLineNum()),
			zap.Float64("query_time", float64(elapsed.Nanoseconds())/1e6),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
			zap.String("agg_type", "gorm"),
		)
	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold != 0 && l.Config.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.Config.SlowThreshold)

		l.Logger.Warn(slowLog,
			zap.String("source", utils.FileWithLineNum()),
			zap.Float64("query_time", float64(elapsed.Nanoseconds())/1e6),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
			zap.String("agg_type", "gorm"),
		)

	case l.Config.LogLevel == gormLogger.Info:
		sql, rows := fc()
		l.Logger.Warn("sql log",
			zap.String("source", utils.FileWithLineNum()),
			zap.Float64("query_time", float64(elapsed.Nanoseconds())/1e6),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
			zap.String("agg_type", "gorm"),
		)
	}
}
