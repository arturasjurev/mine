package manager

import "github.com/sheirys/mine/manager/journal"

type Application struct {
	Journal journal.JournalService
}
