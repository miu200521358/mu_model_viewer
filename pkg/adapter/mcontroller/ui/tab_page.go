//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_model_viewer/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase/minteractor"
)

// NewTabPages は mu_model_viewer のタブページ群を生成する。
func NewTabPages(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialModelPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *minteractor.ModelViewerUsecase) []declarative.TabPage {
	var fileTab *walk.TabPage

	var translator i18n.II18n
	var logger logging.ILogger
	var userConfig config.IUserConfig
	if baseServices != nil {
		translator = baseServices.I18n()
		logger = baseServices.Logger()
		if cfg := baseServices.Config(); cfg != nil {
			userConfig = cfg.UserConfig()
		}
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if viewerUsecase == nil {
		// ユースケース未設定時は依存不足の可能性があるため、空の依存で生成する。
		viewerUsecase = minteractor.NewModelViewerUsecase(minteractor.ModelViewerUsecaseDeps{})
	}

	player := widget.NewMotionPlayer(translator)
	player.SetAudioPlayer(audioPlayer, userConfig)

	materialView := widget.NewMaterialTableView(
		translator,
		i18n.TranslateOrMark(translator, messages.HelpMaterialView),
		func(cw *controller.ControlWindow, indexes []int) {
			if cw == nil {
				return
			}
			cw.SetSelectedMaterialIndexes(0, 0, indexes)
		},
	)

	allMaterialButton := widget.NewMPushButton()
	allMaterialButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelAll))
	allMaterialButton.SetTooltip(i18n.TranslateOrMark(translator, messages.LabelAllTip))
	allMaterialButton.SetMaxSize(declarative.Size{Width: 50})
	allMaterialButton.SetMinSize(declarative.Size{Width: 30})
	allMaterialButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if materialView == nil {
			return
		}
		materialView.SetAllChecked(true)
	})

	invertMaterialButton := widget.NewMPushButton()
	invertMaterialButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelInvert))
	invertMaterialButton.SetTooltip(i18n.TranslateOrMark(translator, messages.LabelInvertTip))
	invertMaterialButton.SetMaxSize(declarative.Size{Width: 50})
	invertMaterialButton.SetMinSize(declarative.Size{Width: 30})
	invertMaterialButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if materialView == nil {
			return
		}
		materialView.InvertChecked()
	})

	var lastModelPath string

	pmxSaveButton := widget.NewMPushButton()
	pmxSaveButton.SetLabel(i18n.TranslateOrMark(translator, messages.LabelPmxSave))
	pmxSaveButton.SetTooltip(i18n.TranslateOrMark(translator, messages.LabelPmxSaveTip))
	pmxSaveButton.SetOnClicked(func(cw *controller.ControlWindow) {
		saveModelAsPmx(viewerUsecase, logger, translator, cw, lastModelPath, 0, 0)
	})

	updatePmxSaveState := func(modelData *model.PmxModel, path string) {
		lastModelPath = ""
		if modelData != nil {
			lastModelPath = modelData.Path()
		}
		if lastModelPath == "" {
			lastModelPath = path
		}
		if pmxSaveButton == nil || pmxSaveButton.PushButton == nil {
			return
		}
		pmxSaveButton.SetEnabled(viewerUsecase.IsPmxConvertiblePath(lastModelPath))
	}

	pmxLoadPicker := widget.NewPmxPmdXLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyPmxHistory,
		i18n.TranslateOrMark(translator, messages.LabelModelFile),
		i18n.TranslateOrMark(translator, messages.LabelModelFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			modelData := loadModel(viewerUsecase, logger, translator, cw, rep, path, materialView, 0, 0)
			updatePmxSaveState(modelData, path)
		},
	)

	vmdLoadPicker := widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelMotionFile),
		i18n.TranslateOrMark(translator, messages.LabelMotionFileTip),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadMotion(viewerUsecase, logger, translator, cw, rep, player, path, 0, 0)
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoadPicker, vmdLoadPicker, materialView, allMaterialButton, invertMaterialButton, pmxSaveButton)

	mWidgets.SetOnLoaded(func() {
		if mWidgets == nil || mWidgets.Window() == nil {
			return
		}
		mWidgets.Window().SetOnEnabledInPlaying(func(playing bool) {
			for _, w := range mWidgets.Widgets {
				w.SetEnabledInPlaying(playing)
			}
		})
		if pmxSaveButton != nil && pmxSaveButton.PushButton != nil {
			pmxSaveButton.SetEnabled(false)
		}
		if initialModelPath != "" {
			pmxLoadPicker.SetPath(initialModelPath)
		}
	})

	fileTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelFile),
		AssignTo: &fileTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					pmxLoadPicker.Widgets(),
					vmdLoadPicker.Widgets(),
					declarative.VSeparator{},
					declarative.Composite{
						Layout: declarative.HBox{},
						Children: []declarative.Widget{
							declarative.TextLabel{Text: i18n.TranslateOrMark(translator, messages.LabelMaterialView)},
							declarative.HSpacer{},
							allMaterialButton.Widgets(),
							invertMaterialButton.Widgets(),
						},
					},
					materialView.Widgets(),
					declarative.VSeparator{},
					pmxSaveButton.Widgets(),
					declarative.VSeparator{},
					player.Widgets(),
				},
			},
		},
	}

	return []declarative.TabPage{fileTabPage}
}

// NewTabPage は単一タブページを生成する。
func NewTabPage(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialModelPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *minteractor.ModelViewerUsecase) declarative.TabPage {
	return NewTabPages(mWidgets, baseServices, initialModelPath, audioPlayer, viewerUsecase)[0]
}

// loadModel はモデル読み込み結果をControlWindowへ反映する。
func loadModel(viewerUsecase *minteractor.ModelViewerUsecase, logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, path string, materialView *widget.MaterialTableView, windowIndex, modelIndex int) *model.PmxModel {
	if cw == nil {
		return nil
	}
	if path == "" {
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	if viewerUsecase == nil {
		logLoadFailed(logger, translator, nil)
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	result, err := viewerUsecase.LoadModel(rep, path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	modelData := (*model.PmxModel)(nil)
	validation := (*minteractor.TextureValidationResult)(nil)
	if result != nil {
		modelData = result.Model
		validation = result.Validation
	}
	if modelData == nil {
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	if materialView != nil {
		logTextureValidationErrors(logger, validation)
		materialView.ResetRows(modelData)
	}
	cw.SetModel(windowIndex, modelIndex, modelData)
	return modelData
}

// loadMotion はモーション読み込み結果をControlWindowへ反映する。
func loadMotion(viewerUsecase *minteractor.ModelViewerUsecase, logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, player *widget.MotionPlayer, path string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if viewerUsecase == nil {
		logLoadFailed(logger, translator, nil)
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	motionResult, err := viewerUsecase.LoadMotion(rep, path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	motionData := (*motion.VmdMotion)(nil)
	maxFrame := motion.Frame(0)
	if motionResult != nil {
		motionData = motionResult.Motion
		maxFrame = motionResult.MaxFrame
	}
	if motionData == nil {
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if player != nil {
		player.Reset(maxFrame)
	}
	cw.SetMotion(windowIndex, modelIndex, motionData)
}

// saveModelAsPmx はXまたはPMDモデルをPMX形式で保存する。
func saveModelAsPmx(viewerUsecase *minteractor.ModelViewerUsecase, logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, modelPath string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	modelData := cw.Model(windowIndex, modelIndex)
	if viewerUsecase == nil {
		logSaveFailed(logger, translator, nil)
		return
	}
	result, err := viewerUsecase.SaveModelAsPmx(minteractor.SaveModelAsPmxRequest{
		ModelPath:              modelPath,
		ModelData:              modelData,
		MissingModelMessage:    i18n.TranslateOrMark(translator, messages.MessageMissingModel),
		InvalidSavePathMessage: i18n.TranslateOrMark(translator, messages.MessageSavePathInvalid),
		SaveOptions:            minteractor.SaveOptions{},
	})
	if err != nil {
		logSaveFailed(logger, translator, err)
		return
	}
	if result == nil || result.OutputPath == "" {
		logSaveFailed(logger, translator, nil)
		return
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	controller.Beep()
	logger.Info(i18n.TranslateOrMark(translator, messages.LogPmxSaveSuccess), filepath.Base(result.OutputPath))
}

// logLoadFailed は読み込み失敗ログを出力する。
func logLoadFailed(logger logging.ILogger, translator i18n.II18n, err error) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageLoadFailed), err)
}

// logSaveFailed は保存失敗ログを出力する。
func logSaveFailed(logger logging.ILogger, translator i18n.II18n, err error) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	logErrorTitle(logger, i18n.TranslateOrMark(translator, messages.MessageSaveFailed), err)
}

// logTextureValidationErrors はテクスチャ検証エラーをログ出力する。
func logTextureValidationErrors(logger logging.ILogger, result *minteractor.TextureValidationResult) {
	if logger == nil || result == nil {
		return
	}
	if len(result.Errors) == 0 {
		return
	}
	for _, err := range result.Errors {
		if err == nil {
			continue
		}
		logger.Warn("テクスチャ検証でエラーが発生しました: %s", err.Error())
	}
}

// logErrorTitle はタイトル付きエラーを出力する。
func logErrorTitle(logger logging.ILogger, title string, err error) {
	if logger == nil {
		return
	}
	if titled, ok := logger.(interface {
		ErrorTitle(title string, err error, msg string, params ...any)
	}); ok {
		titled.ErrorTitle(title, err, "")
		return
	}
	if err == nil {
		logger.Error("%s: %s", title, "")
		return
	}
	logger.Error("%s: %s", title, err.Error())
}
