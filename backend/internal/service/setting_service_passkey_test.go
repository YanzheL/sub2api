//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type passkeySettingsRepoStub struct {
	all map[string]string
}

func (s *passkeySettingsRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
}

func (s *passkeySettingsRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
}

func (s *passkeySettingsRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
}

func (s *passkeySettingsRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
}

func (s *passkeySettingsRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
}

func (s *passkeySettingsRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	out := make(map[string]string, len(s.all))
	for k, v := range s.all {
		out[k] = v
	}
	return out, nil
}

func (s *passkeySettingsRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
}

func TestSettingService_GetAllSettings_PasskeyRPConfig_ExplicitOverridesDerived(t *testing.T) {
	repo := &passkeySettingsRepoStub{
		all: map[string]string{
			SettingKeyPasskeyRPID:           "accounts.example.com",
			SettingKeyPasskeyRPName:         "Example Accounts",
			SettingKeyPasskeyAllowedOrigins: `["https://accounts.example.com"]`,
			SettingKeyFrontendURL:           "https://frontend.example.com",
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	settings, err := svc.GetAllSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, "accounts.example.com", settings.PasskeyRPID)
	require.Equal(t, "Example Accounts", settings.PasskeyRPName)
	require.Equal(t, []string{"https://accounts.example.com"}, settings.PasskeyAllowedOrigins)
}

func TestSettingService_GetAllSettings_PasskeyRPConfig_DerivedFromFrontendURL(t *testing.T) {
	repo := &passkeySettingsRepoStub{
		all: map[string]string{
			SettingKeyFrontendURL: "https://App.Example.com:7443/path",
		},
	}
	svc := NewSettingService(repo, &config.Config{})

	settings, err := svc.GetAllSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, "app.example.com", settings.PasskeyRPID)
	require.Equal(t, "app.example.com", settings.PasskeyRPName)
	require.Equal(t, []string{"https://app.example.com:7443"}, settings.PasskeyAllowedOrigins)
}

func TestSettingService_GetAllSettings_PasskeyRPConfig_DerivedFromConfigFrontendURL(t *testing.T) {
	repo := &passkeySettingsRepoStub{all: map[string]string{}}
	svc := NewSettingService(repo, &config.Config{Server: config.ServerConfig{FrontendURL: "https://config.example.com"}})

	settings, err := svc.GetAllSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, "config.example.com", settings.PasskeyRPID)
	require.Equal(t, "config.example.com", settings.PasskeyRPName)
	require.Equal(t, []string{"https://config.example.com"}, settings.PasskeyAllowedOrigins)
}

func TestSettingService_GetAllSettings_PasskeyRPConfig_UsesLocalhostDefaultsInDebugMode(t *testing.T) {
	repo := &passkeySettingsRepoStub{all: map[string]string{}}
	svc := NewSettingService(repo, &config.Config{Server: config.ServerConfig{Mode: "debug"}})

	settings, err := svc.GetAllSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, "localhost", settings.PasskeyRPID)
	require.Equal(t, "localhost", settings.PasskeyRPName)
	require.Equal(t, passkeyDevelopmentAllowedOrigins, settings.PasskeyAllowedOrigins)
}

func TestSettingService_GetAllSettings_PasskeyRPConfig_DoesNotUseLocalhostDefaultsInReleaseMode(t *testing.T) {
	repo := &passkeySettingsRepoStub{all: map[string]string{}}
	svc := NewSettingService(repo, &config.Config{Server: config.ServerConfig{Mode: "release"}})

	settings, err := svc.GetAllSettings(context.Background())
	require.NoError(t, err)
	require.Equal(t, "", settings.PasskeyRPID)
	require.Equal(t, "Sub2API", settings.PasskeyRPName)
	require.Empty(t, settings.PasskeyAllowedOrigins)
}
