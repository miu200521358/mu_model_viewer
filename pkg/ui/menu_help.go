package ui

import (
	"github.com/miu200521358/walk/pkg/declarative"
)

func GetMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		// declarative.Action{
		// 	Text:        mi18n.T("概要"),
		// 	OnTriggered: func() { mlog.ILT(mi18n.T("概要"), mi18n.T("概要メッセージ")) },
		// },
		declarative.Separator{},
	}
}
