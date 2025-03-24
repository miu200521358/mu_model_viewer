package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/walk/pkg/declarative"
)

func NewMenuItems() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        "MuModelViewer",
			OnTriggered: func() { mlog.ILT("MuModelViewer", "%s", mi18n.T("使い方説明")) },
		},
	}
}
