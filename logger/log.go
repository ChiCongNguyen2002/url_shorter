package logger

import (
  "sync"

  "github.com/google/uuid"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
  "github.com/rs/zerolog/pkgerrors"
)

var (
  loggerInstance *Logger
  mu             sync.RWMutex
)

const (
  KeyServiceName = "service_name"
  KeyLogId       = "log_id"
  KeyFileError   = "file_error"
)

func InitLog(serviceName string) {
  mu.Lock()
  defer mu.Unlock()
  if loggerInstance != nil {
    return
  }

  if serviceName == "" {
    log.Fatal().Msg("service name is empty")
  }

  zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
  lg := log.With().Str(KeyServiceName, serviceName).Logger()
  loggerInstance = &Logger{lg}
}

func GetLogger() *Logger {
  mu.RLock()
  defer mu.RUnlock()
  // handle generate log id default
  uid := uuid.NewString()
  lg := loggerInstance.logger.With().Str(KeyLogId, uid).Logger()
  return &Logger{lg}
}
