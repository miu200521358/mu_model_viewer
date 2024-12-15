//go:build windows
// +build windows

package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/viewer"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
)

var env string

func init() {
	runtime.LockOSThread()

	// システム上の25%の論理プロセッサを使用する
	runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(widget.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

func main() {
	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	mApp := app.NewMApp(appConfig)
	mApp.RunViewerToControlChannel()
	mApp.RunControlToViewerChannel()

	var model *pmx.PmxModel
	if len(os.Args) > 1 {
		modelPath := os.Args[1]
		rep := repository.NewPmxXRepository()
		if isOk, err := rep.CanLoad(modelPath); isOk && err == nil {
			if data, err := rep.Load(modelPath); err == nil {
				model = data.(*pmx.PmxModel)
			} else {
				widget.RaiseError(err)
			}
		} else if err != nil {
			widget.RaiseError(err)
		} else {
			widget.RaiseError(errors.New("unknown error"))
		}
	}

	viewerWindow := viewer.NewViewWindow(
		mApp.ViewerCount(), appConfig, mApp, mApp.ViewerToControlChannel(),
		fmt.Sprintf("%s %s", appConfig.Name, appConfig.Version), nil)
	viewerWindow.GetWindow().SetCloseCallback(func(w *glfw.Window) { mApp.AppState().SetClosed(true) })

	mApp.AddViewWindow(viewerWindow)

	mApp.SetFuncGetModels(func() [][]*pmx.PmxModel {
		return [][]*pmx.PmxModel{{model}}
	})

	motion := vmd.NewVmdMotion("")

	mApp.SetFuncGetMotions(func() [][]*vmd.VmdMotion {
		return [][]*vmd.VmdMotion{{motion}}
	})

	mApp.Center()
	mApp.RunViewer()
}
