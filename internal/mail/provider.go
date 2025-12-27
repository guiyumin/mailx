package mail

// Gmail special folders
const (
	GmailFolderPrefix = "[Gmail]/"
	GmailTrash        = "[Gmail]/Trash"
	GmailAllMail      = "[Gmail]/All Mail"
	GmailDrafts       = "[Gmail]/Drafts"
	GmailSent         = "[Gmail]/Sent Mail"
	GmailStarred      = "[Gmail]/Starred"
	GmailSpam         = "[Gmail]/Spam"
)

// Standard IMAP folders (Yahoo and others)
const (
	INBOX    = "INBOX" // IMAP spec uses uppercase
	Sent     = "Sent"
	Draft    = "Draft"
	Drafts   = "Drafts"
	Trash    = "Trash"
	Spam     = "Spam"
	BulkMail = "Bulk Mail"
	Archive  = "Archive"
	Junk     = "Junk"
)
