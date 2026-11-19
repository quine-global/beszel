package hub

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/auth"
)

func TestSetOIDCProvider(t *testing.T) {
	users := core.NewAuthCollection("users")
	users.OAuth2.Providers = []core.OAuth2ProviderConfig{{Name: "github", ClientId: "gh"}, {Name: auth.NameOIDC, ClientId: "old"}}

	// Unset: nothing changes.
	if err := setOIDCProvider(users); err != nil || users.OAuth2.Enabled || len(users.OAuth2.Providers) != 2 {
		t.Fatalf("unset changed things: %v %+v", err, users.OAuth2)
	}

	secret := filepath.Join(t.TempDir(), "secret")
	os.WriteFile(secret, []byte("s3cret\n"), 0o600)
	t.Setenv("OIDC_AUTH_URL", "https://id.example/auth")
	t.Setenv("OIDC_TOKEN_URL", "https://id.example/token")
	t.Setenv("OIDC_USERINFO_URL", "https://id.example/userinfo")
	t.Setenv("OIDC_CLIENT_ID", "beszel")
	t.Setenv("OIDC_CLIENT_SECRET_FILE", secret)
	t.Setenv("OIDC_DISPLAY_NAME", "Example")
	if err := setOIDCProvider(users); err != nil {
		t.Fatal(err)
	}
	if !users.OAuth2.Enabled || len(users.OAuth2.Providers) != 2 {
		t.Fatalf("got %+v", users.OAuth2)
	}
	p := users.OAuth2.Providers[0]
	if p.Name != auth.NameOIDC || p.ClientId != "beszel" || p.ClientSecret != "s3cret" || p.DisplayName != "Example" || p.PKCE == nil || !*p.PKCE {
		t.Fatalf("got %+v", p)
	}
	if users.OAuth2.Providers[1].Name != "github" {
		t.Fatal("other providers should stay")
	}

	t.Setenv("OIDC_CLIENT_ID", "")
	if err := setOIDCProvider(users); err == nil {
		t.Fatal("missing client id should fail")
	}
}
