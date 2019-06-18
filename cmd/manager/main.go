package main

import (
	"github.com/sheirys/mine/manager"
	"github.com/sheirys/mine/manager/journal"
)

func main() {
	_ = &manager.Application{
		Journal: &journal.JournalFileService{},
	}

}
