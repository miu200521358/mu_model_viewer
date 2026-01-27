//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_model_viewer/pkg/ui"

	"github.com/miu200521358/mlib_go/pkg/infra/app"
	"github.com/miu200521358/mlib_go/pkg/infra/base/err"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/drivers/maudio"
	"github.com/miu200521358/mlib_go/pkg/infra/viewer"
	"github.com/miu200521358/mlib_go/pkg/shared/state"
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
	viewerCount := 1
	initialModelPath := findInitialModelPath(os.Args)

	boot, initErr := app.Init(appFiles, appI18nFiles, nil)
	if initErr != nil {
		if boot != nil {
			err.ShowFatalErrorDialog(boot.AppConfig, initErr)
		} else {
			err.ShowFatalErrorDialog(nil, initErr)
		}
		return
	}
	appConfig := boot.AppConfig
	baseServices := boot.BaseServices
	iconImage := boot.IconImage
	appIcon := boot.AppIcon
	audioPlayer := maudio.NewAudioPlayer()

	sharedState := state.NewSharedState(viewerCount)
	if sharedState == nil {
		err.ShowFatalErrorDialog(appConfig, app.NewSharedStateInitFailed())
		return
	}

	widths, heights, positionXs, positionYs := app.GetCenterSizeAndWidth(appConfig, viewerCount)

	var (
		controlWindow    *controller.ControlWindow
		controlWindowErr error
	)
	viewerManager := viewer.NewViewerManager(sharedState, baseServices)
	if iconImage != nil {
		viewerManager.SetWindowIcon(iconImage)
	}

	go func() {
		defer app.SafeExecute(appConfig, func() {
			widgets := &controller.MWidgets{
				Position: &walk.Point{X: positionXs[0], Y: positionYs[0]},
			}

			controlWindow, controlWindowErr = controller.NewControlWindow(
				sharedState,
				baseServices,
				ui.NewMenuItems(baseServices.I18n(), baseServices.Logger()),
				ui.NewTabPages(widgets, baseServices, initialModelPath, audioPlayer),
				widths[0], heights[0], positionXs[0], positionYs[0], viewerCount,
			)
			if controlWindowErr != nil {
				err.ShowFatalErrorDialog(appConfig, controlWindowErr)
				return
			}
			if appIcon != nil {
				controlWindow.SetIcon(appIcon)
			}

			widgets.SetWindow(controlWindow)
			widgets.OnLoaded()
			controlWindow.Run()
		})
	}()

	if glfwErr := glfw.Init(); glfwErr != nil {
		err.ShowFatalErrorDialog(appConfig, app.NewGlfwInitFailed(glfwErr))
		return
	}

	defer app.SafeExecute(appConfig, func() {
		for n := range viewerCount {
			idx := n + 1
			if addWindowErr := viewerManager.AddWindow(
				fmt.Sprintf("Viewer%d", idx),
				widths[idx], heights[idx], positionXs[idx], positionYs[idx],
			); addWindowErr != nil {
				err.ShowFatalErrorDialog(appConfig, addWindowErr)
				return
			}
		}

		viewerManager.InitOverlay()
		viewerManager.Run()
	})
}

// findInitialModelPath は起動引数からモデルパスを抽出する。
func findInitialModelPath(args []string) string {
	if len(args) <= 1 {
		return ""
	}
	for i := 1; i < len(args); i++ {
		path := strings.TrimSpace(args[i])
		path = strings.Trim(path, "\"")
		path = strings.Trim(path, "'")
		if path == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".pmx", ".pmd", ".x":
			return path
		}
	}
	return ""
}
