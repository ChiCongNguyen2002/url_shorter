package logger

import "context"

const (
  KeyTraceInfo = "trace_info"
)

type TraceInfo struct {
  RequestID string `json:"request_id"`
}

func GetRequestIdByContext(ctx context.Context) *TraceInfo {
  value := ctx.Value(KeyTraceInfo)
  traceInfo, ok := value.(TraceInfo)
  if !ok {
    return nil
  }
  return &traceInfo
}
