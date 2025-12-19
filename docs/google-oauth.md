# How Google OAuth Works

## The Problem OAuth Solves

With App Passwords, you give Maily your actual password (a generated one, but still). If Maily were malicious, it could do anything with your account.

OAuth instead gives apps a **limited, revocable token** - never your actual password.

## The Flow

```
┌─────────┐         ┌─────────┐         ┌─────────┐
│  User   │         │  Maily  │         │ Google  │
└────┬────┘         └────┬────┘         └────┬────┘
     │                   │                   │
     │  1. Run login     │                   │
     │──────────────────>│                   │
     │                   │                   │
     │  2. Opens browser │                   │
     │<──────────────────│                   │
     │                   │                   │
     │  3. "Allow Maily to access Gmail?"    │
     │──────────────────────────────────────>│
     │                   │                   │
     │  4. User clicks "Allow"               │
     │──────────────────────────────────────>│
     │                   │                   │
     │  5. Redirect to localhost with code   │
     │<──────────────────────────────────────│
     │                   │                   │
     │                   │  6. Exchange code │
     │                   │     for tokens    │
     │                   │──────────────────>│
     │                   │                   │
     │                   │  7. Access token  │
     │                   │     + Refresh token
     │                   │<──────────────────│
     │                   │                   │
     │  8. Done! Token   │                   │
     │     stored locally│                   │
     │<──────────────────│                   │
```

## The Key Parts

**1. Client ID & Secret** - You register Maily with Google Cloud. Google gives you credentials that identify your app.

**2. Authorization URL** - Maily opens a browser to:
```
https://accounts.google.com/o/oauth2/auth?
  client_id=YOUR_CLIENT_ID&
  redirect_uri=http://localhost:8080/callback&
  scope=https://mail.google.com/&
  response_type=code
```

**3. User Consent** - Google shows "Maily wants to access your Gmail. Allow?"

**4. Authorization Code** - User clicks Allow, Google redirects to `localhost:8080/callback?code=ABC123`

**5. Token Exchange** - Maily exchanges that code for:
- **Access Token** - Short-lived (~1 hour), used to actually access Gmail
- **Refresh Token** - Long-lived, used to get new access tokens

**6. Using the Token** - Instead of username/password for IMAP, you use:
```
Username: user@gmail.com
Password: <access_token>
Auth method: XOAUTH2
```

## Why It's Better

| App Password              | OAuth                                |
|---------------------------|--------------------------------------|
| Never expires             | Tokens expire, auto-refresh          |
| Full account access       | Scoped permissions (just Gmail)      |
| User must manually revoke | User can revoke in Google settings   |
| No audit trail            | Google logs which apps accessed what |

## What Maily Would Need

1. A Google Cloud project with OAuth credentials
2. Code to open browser and run a tiny localhost server for the callback
3. Code to exchange the code for tokens
4. Store refresh token in `~/.config/maily/`
5. Use XOAUTH2 authentication with go-imap instead of plain password

---

# Implementation Plan

## User Experience

```
$ maily login gmail --oauth

  ██████████████████████████████████
  ██████████████████████████████████
  ████ ▄▄▄▄▄ █ ▄▄ █▄█▄█ ▄▄▄▄▄ ████
  ████ █   █ █ ▀▀██▀▄ █ █   █ ████
  ████ █▄▄▄█ █▀▄▀▀▀▄▄▀█ █▄▄▄█ ████
  ████▄▄▄▄▄▄▄█▄▀▄█▄▀▄█▄▄▄▄▄▄▄████
  ██████████████████████████████████

  Scan QR code or open this URL:
  https://accounts.google.com/o/oauth2/auth?...

  Waiting for authorization...
  (or press 'c' to enter code manually)

  ✓ Successfully logged in as user@gmail.com
```

## Features

1. **QR code in terminal** - Scan with phone to open auth URL
2. **Localhost callback** - Automatic token capture (primary)
3. **Manual code entry** - Fallback if localhost doesn't work

## Files to Create

### `internal/auth/oauth.go`
OAuth flow implementation:
- `StartOAuthFlow()` - Main entry point
- `startLocalServer()` - Localhost callback server on random port
- `generateAuthURL()` - Build Google OAuth URL with PKCE
- `exchangeCodeForTokens()` - Token exchange
- `RefreshAccessToken()` - Token refresh logic

### `internal/auth/oauth_config.go`
```go
var OAuthConfig = &oauth2.Config{
    ClientID:     "EMBEDDED_CLIENT_ID",
    ClientSecret: "EMBEDDED_CLIENT_SECRET",
    Scopes:       []string{"https://mail.google.com/"},
    Endpoint:     google.Endpoint,
}
```

### `internal/ui/oauth_login.go`
TUI component with states:
- `oauthStateWaiting` - Display QR code + URL
- `oauthStateManualEntry` - Text input for code
- `oauthStateExchanging` - Spinner during exchange
- `oauthStateSuccess` / `oauthStateError`

## Files to Modify

### `internal/auth/credentials.go`
```go
type Credentials struct {
    Email        string `yaml:"email"`
    Password     string `yaml:"password,omitempty"`      // App password auth
    AccessToken  string `yaml:"access_token,omitempty"`  // OAuth
    RefreshToken string `yaml:"refresh_token,omitempty"` // OAuth
    TokenExpiry  int64  `yaml:"token_expiry,omitempty"`  // Unix timestamp
    AuthMethod   string `yaml:"auth_method"`             // "password" or "oauth"
    IMAPHost     string `yaml:"imap_host"`
    IMAPPort     int    `yaml:"imap_port"`
    SMTPHost     string `yaml:"smtp_host"`
    SMTPPort     int    `yaml:"smtp_port"`
}
```

### `internal/cli/login.go`
Add `--oauth` flag to route to OAuth flow.

### `internal/gmail/imap.go`
- Detect OAuth vs password auth
- Use XOAUTH2 SASL mechanism
- Auto-refresh expired tokens before connecting

### `internal/gmail/smtp.go`
Support XOAUTH2 for sending emails.

## Dependencies

```bash
go get github.com/mdp/qrterminal/v3  # QR code in terminal
go get golang.org/x/oauth2           # OAuth 2.0 client
go get golang.org/x/oauth2/google    # Google endpoints
```

## OAuth URL Structure

```
https://accounts.google.com/o/oauth2/v2/auth?
  client_id=CLIENT_ID&
  redirect_uri=http://localhost:PORT/callback&
  response_type=code&
  scope=https://mail.google.com/&
  state=RANDOM_STATE&
  code_challenge=PKCE_CHALLENGE&
  code_challenge_method=S256&
  access_type=offline&
  prompt=consent
```

## Security

- **PKCE** (Proof Key for Code Exchange) for authorization flow
- **State parameter** to prevent CSRF attacks
- **File permissions** 0600 on accounts.yml (already implemented)
- **Auto-refresh** tokens before expiry

## Setup (One-time)

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create project "Maily"
3. Enable Gmail API
4. Configure OAuth consent screen (External)
5. Create OAuth 2.0 credentials (Desktop app type)
6. Embed Client ID and Secret in `oauth_config.go`

## Notes

- Unverified apps show Google warning (OK for personal use)
- For public distribution, need Google app verification
- App passwords remain supported as alternative
