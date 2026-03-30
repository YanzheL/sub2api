package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasskeyFriendlyNameResolver_ResolvePrecedence(t *testing.T) {
	now := time.Date(2026, 3, 30, 9, 30, 0, 0, time.UTC)
	metadataOnlyAAGUID := "11111111-2222-4333-8444-555555555555"

	resolver := newPasskeyFriendlyNameResolver(NewStaticPasskeyAAGUIDMetadataCache(map[string]string{
		passkeyBitwardenAAGUID: "Bitwarden (Metadata)",
		metadataOnlyAAGUID:     "FIDO Metadata Key",
	}))

	require.Equal(t, "Personal Key", resolver.Resolve(t.Context(), "  Personal Key  ", passkeyBitwardenAAGUID, now))
	require.Equal(t, "Bitwarden Passkey", resolver.Resolve(t.Context(), "", passkeyBitwardenAAGUID, now))
	require.Equal(t, "FIDO Metadata Key", resolver.Resolve(t.Context(), "", metadataOnlyAAGUID, now))
	require.Equal(t, passkeyFriendlyName("", now), resolver.Resolve(t.Context(), "", "unmapped-aaguid", now))
}

func TestNewPasskeyAAGUIDMetadataCacheFromJSON_SupportsMapPayload(t *testing.T) {
	payload := []byte(`{"d548826e-79b4-db40-a3d8-11116f7e8349":"Bitwarden Metadata"}`)

	cache, err := NewPasskeyAAGUIDMetadataCacheFromJSON(payload)
	require.NoError(t, err)
	require.NotNil(t, cache)

	name, ok := cache.LookupFriendlyNameByAAGUID(t.Context(), passkeyBitwardenAAGUID)
	require.True(t, ok)
	require.Equal(t, "Bitwarden Metadata", name)
}

func TestNewPasskeyAAGUIDMetadataCacheFromJSON_SupportsListPayload(t *testing.T) {
	payload := []byte(`[{"aaguid":"11111111-2222-4333-8444-555555555555","name":"Metadata Hardware Key"}]`)

	cache, err := NewPasskeyAAGUIDMetadataCacheFromJSON(payload)
	require.NoError(t, err)
	require.NotNil(t, cache)

	name, ok := cache.LookupFriendlyNameByAAGUID(t.Context(), "11111111-2222-4333-8444-555555555555")
	require.True(t, ok)
	require.Equal(t, "Metadata Hardware Key", name)
}

func TestLoadOptionalPasskeyAAGUIDMetadataCacheFromEnv_LoadsCache(t *testing.T) {
	dir := t.TempDir()
	metadataPath := filepath.Join(dir, "aaguid-metadata.json")
	err := os.WriteFile(metadataPath, []byte(`{"11111111-2222-4333-8444-555555555555":"Env Metadata Key"}`), 0o600)
	require.NoError(t, err)

	t.Setenv(passkeyAAGUIDMetadataCachePathEnv, metadataPath)

	cache := loadOptionalPasskeyAAGUIDMetadataCacheFromEnv()
	require.NotNil(t, cache)

	name, ok := cache.LookupFriendlyNameByAAGUID(t.Context(), "11111111-2222-4333-8444-555555555555")
	require.True(t, ok)
	require.Equal(t, "Env Metadata Key", name)
}
