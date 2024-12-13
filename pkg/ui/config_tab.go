package ui

import (
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func newConfigTab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.ConfigTab = widget.NewMTabPage(mi18n.T("設定"))
	controlWindow.AddTabPage(toolState.ConfigTab.TabPage)

	toolState.ConfigTab.SetLayout(walk.NewVBoxLayout())
	composite := &declarative.Composite{
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text: mi18n.T("設定説明"),
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					// サイジングセット設定読み込みボタン
					declarative.PushButton{
						Text: mi18n.T("設定保存"),
						OnClicked: func() {
							// toolState.loadSizingSet()
						},
					},
				},
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					// サイジングセット設定読み込みボタン
					declarative.PushButton{
						Text: mi18n.T("設定保存"),
						OnClicked: func() {
							// toolState.loadSizingSet()
						},
					},
				},
			},
		},
	}

	if err := composite.Create(declarative.NewBuilder(toolState.ConfigTab)); err != nil {
		widget.RaiseError(err)
	}
}
