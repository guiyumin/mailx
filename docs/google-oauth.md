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
