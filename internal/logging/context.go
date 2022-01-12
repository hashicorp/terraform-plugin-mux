package logging

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

// InitContext creates SDK logger contexts.
func InitContext(ctx context.Context) context.Context {
	ctx = tfsdklog.NewSubsystem(ctx, SubsystemMux, tfsdklog.WithLevelFromEnv(EnvTfLogSdkMux))

	return ctx
}

// Tfprotov5ProviderServerContext injects the chosen provider Go type
func Tfprotov5ProviderServerContext(ctx context.Context, p tfprotov5.ProviderServer) context.Context {
	providerType := fmt.Sprintf("%T", p)
	ctx = tflog.With(ctx, KeyTfMuxProvider, providerType)
	ctx = tfsdklog.With(ctx, KeyTfMuxProvider, providerType)
	ctx = tfsdklog.SubsystemWith(ctx, SubsystemMux, KeyTfMuxProvider, providerType)

	return ctx
}
