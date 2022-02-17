package helpers

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func CustomLogs(message string, ctx context.Context, enabled bool) {
	customLogs := log.FromContext(ctx)
	customLogs.Info("Info:[ " + message + " ]\n")
}
