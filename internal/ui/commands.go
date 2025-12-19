package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"maily/internal/gmail"
)

func (a App) initClient() tea.Cmd {
	account := a.currentAccount()
	if account == nil {
		return func() tea.Msg {
			return errorMsg{err: fmt.Errorf("no account configured")}
		}
	}
	creds := &account.Credentials
	return func() tea.Msg {
		client, err := gmail.NewIMAPClient(creds)
		if err != nil {
			return errorMsg{err: err}
		}
		return clientReadyMsg{imap: client}
	}
}

func (a *App) loadEmails() tea.Cmd {
	return func() tea.Msg {
		emails, err := a.imap.FetchMessages("INBOX", a.emailLimit)
		if err != nil {
			return errorMsg{err: err}
		}
		return emailsLoadedMsg{emails: emails}
	}
}

func (a *App) executeSearch(query string) tea.Cmd {
	return func() tea.Msg {
		emails, err := a.imap.SearchMessages("INBOX", query)
		if err != nil {
			return errorMsg{err: err}
		}
		return appSearchResultsMsg{emails: emails, query: query}
	}
}
