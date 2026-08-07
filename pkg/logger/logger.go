package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger ตั้งค่าและสร้าง Zap Logger ตาม Environment (development / production)
func InitLogger(env string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	Log, err = config.Build(zap.AddCallerSkip(1)) // Skip 1 level เพื่อแสดง caller ในไฟล์ที่ถูกเรียกจริง
	if err != nil {
		os.Exit(1)
	}

	zap.ReplaceGlobals(Log)
}

// ensureInitialized ช่วยป้องกัน Nil Pointer Exception ในกรณีที่ลืมเรียก InitLogger() ก่อนใช้งาน
func getLogger() *zap.Logger {
	if Log == nil {
		InitLogger("development") // Default เป็น development หากยังไม่ได้ init
	}
	return Log
}

// --- Logging Methods ---

func Info(message string, fields ...zap.Field) {
	getLogger().Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	getLogger().Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	getLogger().Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	getLogger().Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	getLogger().Fatal(message, fields...)
}

// Sync เคลียร์ Log ที่ค้างใน Buffer
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// --- Zap Field Helpers (เพิ่มส่วนนี้เพื่อให้ Repository เรียกใช้งานได้สะดวก) ---

func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
