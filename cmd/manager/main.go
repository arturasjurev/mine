package main

import (
	"github.com/sheirys/mine/manager"
	"github.com/sheirys/mine/manager/journal"
)

func main() {
	app := &manager.Application{
		Journal: journal.JournalFileService{},
	}
}
