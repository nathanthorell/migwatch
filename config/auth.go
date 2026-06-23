package config

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nathanthorell/migwatch/model"
)

const postgresTokenScope = "https://ossrdbms-aad.database.windows.net/.default"

func WrapAuthError(err error, conn model.Connection) error {
	msg := err.Error()

	isAuthErr := strings.Contains(msg, "Login failed") ||
		strings.Contains(msg, "login error") ||
		strings.Contains(msg, "password authentication failed") ||
		strings.Contains(msg, "no pg_hba.conf entry")

	if !isAuthErr {
		return err
	}

	switch conn.AuthMethod {
	case "ActiveDirectoryAzCli":
		return fmt.Errorf("authentication failed: az login token missing or expired — run `az login`")
	case "ActiveDirectoryInteractive":
		return fmt.Errorf("authentication failed: interactive login did not complete")
	case "ActiveDirectoryDefault":
		return fmt.Errorf("authentication failed: no valid credential found in default chain (az login, env vars, managed identity)")
	default:
		return fmt.Errorf("authentication failed: invalid username or password")
	}
}

// ResolveConnection injects an access token into the DSN when the driver and auth method require one.
func ResolveConnection(ctx context.Context, conn model.Connection) (model.Connection, error) {
	if conn.Driver != model.DriverPostgres || conn.AuthMethod == "" {
		return conn, nil
	}

	token, err := fetchPostgresToken(ctx, conn.AuthMethod)
	if err != nil {
		return conn, err
	}

	injected, err := injectPassword(conn.DSN, token)
	if err != nil {
		return conn, err
	}

	conn.DSN = injected
	return conn, nil
}

func fetchPostgresToken(ctx context.Context, authMethod string) (string, error) {
	var cred interface {
		GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error)
	}

	switch authMethod {
	case "ActiveDirectoryDefault":
		c, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return "", fmt.Errorf("create DefaultAzureCredential: %w", err)
		}
		cred = c
	case "ActiveDirectoryAzCli":
		c, err := azidentity.NewAzureCLICredential(nil)
		if err != nil {
			return "", fmt.Errorf("create AzureCliCredential: %w", err)
		}
		cred = c
	default:
		return "", fmt.Errorf("unsupported postgres fedauth method %q", authMethod)
	}

	tk, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{postgresTokenScope}})
	if err != nil {
		return "", fmt.Errorf("get token (%s): %w", authMethod, err)
	}
	return tk.Token, nil
}

// injectPassword sets the password in a postgres DSN and removes the fedauth param.
func injectPassword(dsn, password string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("parse DSN: %w", err)
	}

	u.User = url.UserPassword(u.User.Username(), password)

	q := u.Query()
	q.Del("fedauth")
	u.RawQuery = q.Encode()

	return u.String(), nil
}
