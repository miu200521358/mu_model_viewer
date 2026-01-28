//go:build windows
// +build windows

// 指示: miu200521358
package main

import (
	"embed"
	"os"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_model_viewer/pkg/adapter/mcontroller/ui"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase/minteractor"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_motion"
	"github.com/miu200521358/mlib_go/pkg/infra/app"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
)

// init はOSスレッド固定とコンソール登録を行う。
func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(controller.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

// main は mu_model_viewer を起動する。
func main() {
	initialModelPath := app.FindInitialPath(os.Args, ".pmx", ".pmd", ".x")

	app.Run(app.RunOptions{
		ViewerCount: 1,
		AppFiles:    appFiles,
		I18nFiles:   appI18nFiles,
		BuildMenuItems: func(baseServices base.IBaseServices) []declarative.MenuItem {
			return ui.NewMenuItems(baseServices.I18n(), baseServices.Logger())
		},
		BuildTabPages: func(widgets *controller.MWidgets, baseServices base.IBaseServices, audioPlayer audio_api.IAudioPlayer) []declarative.TabPage {
			viewerUsecase := minteractor.NewModelViewerUsecase(minteractor.ModelViewerUsecaseDeps{
				ModelReader:      io_model.NewModelRepository(),
				MotionReader:     io_motion.NewVmdVpdRepository(),
				ModelWriter:      io_model.NewModelRepository(),
				TextureValidator: mfile.NewTextureValidator(),
				PathService:      mfile.NewPathService(),
			})
			return ui.NewTabPages(widgets, baseServices, initialModelPath, audioPlayer, viewerUsecase)
		},
	})
}
