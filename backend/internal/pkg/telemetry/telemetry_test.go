package telemetry

import (
    "testing"
    "context"
)

func TestTelemetry(t *testing.T) {
    shutdown, err := InitTelemetry(context.Background(), "test-proj", "test-srv")
    if err == nil && shutdown != nil {
        shutdown(context.Background())
    }
}
