package ui

import (
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
)

type ToolState struct {
	App           *app.MApp
	ControlWindow *controller.ControlWindow
	ConfigTab     *widget.MTabPage
}

func NewToolState(app *app.MApp, controlWindow *controller.ControlWindow) *ToolState {

	toolState := &ToolState{
		App:           app,
		ControlWindow: controlWindow,
	}

	newConfigTab(controlWindow, toolState)

	return toolState
}
