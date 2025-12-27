# Refactor PRD: Provider Constants and Search Cleanup

## Goals

1. Extract magic strings to constants
2. Reduce scattered conditional logic
3. Unify connection/login logic in search.go (minimal abstraction)

## Non-Goals

- No Provider interface or abstract factory patterns
- No YAML configuration files for providers
- No Strategy pattern for search

---

## 1. Extract Magic Strings to Constants

**Create `internal/mail/provider.go`** with all provider-related constants:

```go
package mail

// Provider identifiers
const (
    ProviderGmail = "gmail"
    ProviderYahoo = "yahoo"
)

// Gmail IMAP/SMTP hosts
const (
    GmailIMAPHost = "imap.gmail.com"
    GmailSMTPHost = "smtp.gmail.com"
)

// Yahoo IMAP/SMTP hosts
const (
    YahooIMAPHost = "imap.mail.yahoo.com"
    YahooSMTPHost = "smtp.mail.yahoo.com"
)

// Standard ports
const (
    IMAPPort = 993
    SMTPPort = 587
)

// Gmail special folders
const (
    GmailFolderPrefix  = "[Gmail]/"
    GmailTrash         = "[Gmail]/Trash"
    GmailAllMail       = "[Gmail]/All Mail"
    GmailDrafts        = "[Gmail]/Drafts"
    GmailSent          = "[Gmail]/Sent Mail"
    GmailStarred       = "[Gmail]/Starred"
    GmailSpam          = "[Gmail]/Spam"
)

// IsGmailHost checks if the host is Gmail
func IsGmailHost(host string) bool {
    return host == GmailIMAPHost
}
```

**Files to update:**
- `internal/auth/credentials.go` - Use constants in `GmailCredentials()` and `YahooCredentials()`
- `internal/mail/imap.go` - Replace `"[Gmail]/Trash"` etc. with constants
- `internal/mail/search.go` - Replace `strings.Contains(creds.IMAPHost, "gmail")` with `IsGmailHost()`
- `internal/ui/components/labelpicker.go` - Use constants for folder mappings

---

## 2. Reduce Scattered Conditional Logic

**Add `Provider` field to `Credentials` struct** (credentials.go):

```go
type Credentials struct {
    Email    string `yaml:"email"`
    Password string `yaml:"password"`
    IMAPHost string `yaml:"imap_host"`
    IMAPPort int    `yaml:"imap_port"`
    SMTPHost string `yaml:"smtp_host"`
    SMTPPort int    `yaml:"smtp_port"`
    Provider string `yaml:"provider"` // "gmail", "yahoo", etc.
}
```

Update factory functions to set this field:

```go
func GmailCredentials(email, password string) Credentials {
    return Credentials{
        Email:    email,
        Password: password,
        IMAPHost: mail.GmailIMAPHost,
        IMAPPort: mail.IMAPPort,
        SMTPHost: mail.GmailSMTPHost,
        SMTPPort: mail.SMTPPort,
        Provider: mail.ProviderGmail,
    }
}
```

**Benefit:** No more host string matching. Check `creds.Provider == mail.ProviderGmail` instead.

---

## 3. Unify Search Connection Logic

**Current state:** `GmailSearch()` and `StandardSearch()` are 95% identical. Only the search command differs.

**Proposed change:** Single function with search type parameter.

```go
// searchType determines which IMAP search command to use
type searchType int

const (
    searchTypeText   searchType = iota // Standard IMAP TEXT search
    searchTypeGmailRaw                  // Gmail X-GM-RAW extension
)

func doSearch(creds *auth.Credentials, mailbox, query string, stype searchType) ([]imap.UID, error) {
    addr := fmt.Sprintf("%s:%d", creds.IMAPHost, creds.IMAPPort)

    conn, err := tls.Dial("tcp", addr, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)

    // Read greeting
    if _, err := reader.ReadString('\n'); err != nil {
        return nil, fmt.Errorf("failed to read greeting: %w", err)
    }

    // Login
    loginCmd := fmt.Sprintf("a1 LOGIN %s %s\r\n", quoteString(creds.Email), quoteString(creds.Password))
    if _, err := conn.Write([]byte(loginCmd)); err != nil {
        return nil, fmt.Errorf("failed to send login: %w", err)
    }
    if err := readUntilOK(reader, "a1"); err != nil {
        return nil, fmt.Errorf("login failed: %w", err)
    }

    // Select mailbox
    selectCmd := fmt.Sprintf("a2 SELECT %s\r\n", quoteString(mailbox))
    if _, err := conn.Write([]byte(selectCmd)); err != nil {
        return nil, fmt.Errorf("failed to send select: %w", err)
    }
    if err := readUntilOK(reader, "a2"); err != nil {
        return nil, fmt.Errorf("select failed: %w", err)
    }

    // Build search command based on type
    var searchCmd string
    switch stype {
    case searchTypeGmailRaw:
        searchCmd = fmt.Sprintf("a3 UID SEARCH X-GM-RAW %s\r\n", quoteString(query))
    default:
        searchCmd = fmt.Sprintf("a3 UID SEARCH TEXT %s\r\n", quoteString(query))
    }

    if _, err := conn.Write([]byte(searchCmd)); err != nil {
        return nil, fmt.Errorf("failed to send search: %w", err)
    }

    uids, err := readSearchResponse(reader, "a3")
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }

    conn.Write([]byte("a4 LOGOUT\r\n"))
    return uids, nil
}

// Search performs a search using the appropriate method for the provider
func Search(creds *auth.Credentials, mailbox, query string) ([]imap.UID, error) {
    if creds.Provider == ProviderGmail {
        return doSearch(creds, mailbox, query, searchTypeGmailRaw)
    }
    return doSearch(creds, mailbox, query, searchTypeText)
}
```

---

## Decisions

1. **No backward compat wrappers** - Remove `GmailSearch()` and `StandardSearch()` entirely
2. **Self-describing credentials** - Add `Provider` to `Credentials` (some duplication with `Account.Provider` is acceptable)

---

## Summary of Changes

| File | Change |
|------|--------|
| `internal/mail/provider.go` | NEW - Constants for hosts, folders, provider IDs |
| `internal/auth/credentials.go` | Add `Provider` field, use constants |
| `internal/mail/search.go` | Unify to single `doSearch()`, remove old functions |
| `internal/mail/imap.go` | Use folder constants |
| `internal/ui/components/labelpicker.go` | Use folder constants |
