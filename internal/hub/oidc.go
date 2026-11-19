package hub

import (
	"errors"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/auth"
)

// setOIDCProvider configures an OpenID Connect sign-in provider (Keycloak,
// Authentik, ...) for the users collection from the environment, so it can be
// set up without the PocketBase dashboard:
//
//	OIDC_AUTH_URL, OIDC_TOKEN_URL, OIDC_USERINFO_URL  the provider's endpoints
//	OIDC_CLIENT_ID                                     the client
//	OIDC_CLIENT_SECRET or OIDC_CLIENT_SECRET_FILE      its secret (the file wins)
//	OIDC_DISPLAY_NAME                                  the button's label
//
// With OIDC_AUTH_URL unset nothing changes, and a provider set up in the
// dashboard stays. Combine with USER_CREATION=true to let the provider's
// accounts sign in without being created in Beszel first.
func setOIDCProvider(users *core.Collection) error {
	authURL, _ := GetEnv("OIDC_AUTH_URL")
	if authURL == "" {
		return nil
	}
	tokenURL, _ := GetEnv("OIDC_TOKEN_URL")
	userInfoURL, _ := GetEnv("OIDC_USERINFO_URL")
	clientID, _ := GetEnv("OIDC_CLIENT_ID")
	secret, _ := GetEnv("OIDC_CLIENT_SECRET")
	if file, _ := GetEnv("OIDC_CLIENT_SECRET_FILE"); file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		secret = strings.TrimSpace(string(b))
	}
	if tokenURL == "" || userInfoURL == "" || clientID == "" || secret == "" {
		return errors.New("OIDC_AUTH_URL is set, but not all of OIDC_TOKEN_URL, OIDC_USERINFO_URL, OIDC_CLIENT_ID and OIDC_CLIENT_SECRET(_FILE)")
	}
	name, _ := GetEnv("OIDC_DISPLAY_NAME")
	if name == "" {
		name = "OpenID Connect"
	}

	pkce := true
	provider := core.OAuth2ProviderConfig{
		Name:         auth.NameOIDC,
		ClientId:     clientID,
		ClientSecret: secret,
		AuthURL:      authURL,
		TokenURL:     tokenURL,
		UserInfoURL:  userInfoURL,
		DisplayName:  name,
		PKCE:         &pkce,
	}
	providers := []core.OAuth2ProviderConfig{provider}
	for _, p := range users.OAuth2.Providers {
		if p.Name != auth.NameOIDC {
			providers = append(providers, p)
		}
	}
	users.OAuth2.Providers = providers
	users.OAuth2.Enabled = true
	return nil
}
