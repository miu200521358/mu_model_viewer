//go:build windows
// +build windows

package ui

import (
	"path/filepath"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/adapter/io_model"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/infra/file/mfile"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mlib_go/pkg/shared/base/merr"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

// overrideBoneInserter は不足ボーンの補完を行うI/F。
type overrideBoneInserter interface {
	InsertShortageOverrideBones() error
}

const (
	repositoryNotConfiguredErrorID = "95504"
	modelNotLoadedErrorID          = "95505"
	savePathInvalidErrorID         = "95506"
)

// NewTabPages は mu_model_viewer のタブページ群を生成する。
func NewTabPages(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialModelPath string, audioPlayer audio_api.IAudioPlayer) []declarative.TabPage {
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

	player := widget.NewMotionPlayer(translator)
	player.SetAudioPlayer(audioPlayer, userConfig)

	materialView := widget.NewMaterialTableView(
		translator,
		translate(translator, "材質ビュー説明"),
		func(cw *controller.ControlWindow, indexes []int) {
			if cw == nil {
				return
			}
			cw.SetSelectedMaterialIndexes(0, 0, indexes)
		},
	)

	allMaterialButton := widget.NewMPushButton()
	allMaterialButton.SetLabel(translate(translator, "全"))
	allMaterialButton.SetTooltip(translate(translator, "全ボタン説明"))
	allMaterialButton.SetMaxSize(declarative.Size{Width: 50})
	allMaterialButton.SetMinSize(declarative.Size{Width: 30})
	allMaterialButton.SetOnClicked(func(cw *controller.ControlWindow) {
		if materialView == nil {
			return
		}
		materialView.SetAllChecked(true)
	})

	invertMaterialButton := widget.NewMPushButton()
	invertMaterialButton.SetLabel(translate(translator, "反"))
	invertMaterialButton.SetTooltip(translate(translator, "反ボタン説明"))
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
	pmxSaveButton.SetLabel(translate(translator, "PMX保存"))
	pmxSaveButton.SetTooltip(translate(translator, "PMX保存説明"))
	pmxSaveButton.SetOnClicked(func(cw *controller.ControlWindow) {
		saveModelAsPmx(logger, translator, cw, lastModelPath, 0, 0)
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
		pmxSaveButton.SetEnabled(isPmxConvertiblePath(lastModelPath))
	}

	pmxLoadPicker := widget.NewPmxXLoadFilePicker(
		userConfig,
		translator,
		"pmx",
		translate(translator, "モデルファイル"),
		translate(translator, "モデルファイルを選択してください"),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			modelData := loadModel(logger, translator, cw, rep, path, materialView, 0, 0)
			updatePmxSaveState(modelData, path)
		},
	)

	vmdLoadPicker := widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		"vmd",
		translate(translator, "モーションファイル"),
		translate(translator, "モーションファイルを選択してください"),
		func(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
			loadMotion(logger, translator, cw, rep, player, path, 0, 0)
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
		Title:    translate(translator, "ファイル"),
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
							declarative.TextLabel{Text: translate(translator, "材質ビュー")},
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
func NewTabPage(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialModelPath string, audioPlayer audio_api.IAudioPlayer) declarative.TabPage {
	return NewTabPages(mWidgets, baseServices, initialModelPath, audioPlayer)[0]
}

// loadModel はモデル読み込み結果をControlWindowへ反映する。
func loadModel(logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, path string, materialView *widget.MaterialTableView, windowIndex, modelIndex int) *model.PmxModel {
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
	if rep == nil {
		logLoadFailed(logger, translator, newRepositoryNotConfiguredError(translate(translator, "モデル読み込みリポジトリがありません")))
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	data, err := rep.Load(path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		logLoadFailed(logger, translator, io_common.NewIoFormatNotSupported(translate(translator, "モデル形式が不正です"), nil))
		if materialView != nil {
			materialView.ResetRows(nil)
		}
		cw.SetModel(windowIndex, modelIndex, nil)
		return nil
	}
	if modelData.Bones != nil {
		if inserter, ok := any(modelData.Bones).(overrideBoneInserter); ok {
			// システム用ボーンの追加に失敗しても何もメッセージは出さない
			inserter.InsertShortageOverrideBones()
		}
	}
	if materialView != nil {
		validateModelTextures(modelData)
		materialView.ResetRows(modelData)
	}
	cw.SetModel(windowIndex, modelIndex, modelData)
	return modelData
}

// validateModelTextures はモデルのテクスチャ有効性を検証する。
func validateModelTextures(modelData *model.PmxModel) {
	if modelData == nil || modelData.Textures == nil {
		return
	}

	baseDir := filepath.Dir(modelData.Path())
	for _, texture := range modelData.Textures.Values() {
		if texture == nil {
			continue
		}
		name := texture.Name()
		if name == "" {
			texture.SetValid(false)
			continue
		}
		texturePath := name
		if !filepath.IsAbs(texturePath) {
			texturePath = filepath.Join(baseDir, texturePath)
		}
		exists, err := mfile.ExistsFile(texturePath)
		if err != nil || !exists {
			texture.SetValid(false)
			continue
		}
		if _, err := mfile.LoadImage(texturePath); err != nil {
			texture.SetValid(false)
			continue
		}
		texture.SetValid(true)
	}
}

// loadMotion はモーション読み込み結果をControlWindowへ反映する。
func loadMotion(logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, rep io_common.IFileReader, player *widget.MotionPlayer, path string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if rep == nil {
		logLoadFailed(logger, translator, newRepositoryNotConfiguredError(translate(translator, "モーション読み込みリポジトリがありません")))
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	data, err := rep.Load(path)
	if err != nil {
		logLoadFailed(logger, translator, err)
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		logLoadFailed(logger, translator, io_common.NewIoFormatNotSupported(translate(translator, "モーション形式が不正です"), nil))
		cw.SetMotion(windowIndex, modelIndex, nil)
		return
	}
	if player != nil {
		player.Reset(motionData.MaxFrame())
	}
	cw.SetMotion(windowIndex, modelIndex, motionData)
}

// saveModelAsPmx はXまたはPMDモデルをPMX形式で保存する。
func saveModelAsPmx(logger logging.ILogger, translator i18n.II18n, cw *controller.ControlWindow, modelPath string, windowIndex, modelIndex int) {
	if cw == nil {
		return
	}
	if !isPmxConvertiblePath(modelPath) {
		logSaveFailed(logger, translator, newModelNotLoadedError(translate(translator, "XまたはPMDファイルが読み込まれていません")))
		return
	}
	modelData := cw.Model(windowIndex, modelIndex)
	if modelData == nil {
		logSaveFailed(logger, translator, newModelNotLoadedError(translate(translator, "XまたはPMDファイルが読み込まれていません")))
		return
	}
	outputPath := buildPmxOutputPath(modelPath)
	if outputPath == "" || !mfile.CanSave(outputPath) {
		logSaveFailed(logger, translator, newSavePathInvalidError(translate(translator, "保存先パスが不正です")))
		return
	}
	if err := io_model.NewModelRepository().Save(outputPath, modelData, io_common.SaveOptions{}); err != nil {
		logSaveFailed(logger, translator, err)
		return
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	logger.Info("PMX保存完了", filepath.Base(outputPath))
}

// newRepositoryNotConfiguredError はリポジトリ未設定エラーを生成する。
func newRepositoryNotConfiguredError(message string) error {
	return merr.NewCommonError(repositoryNotConfiguredErrorID, merr.ErrorKindInternal, message, nil)
}

// newModelNotLoadedError はモデル未読込エラーを生成する。
func newModelNotLoadedError(message string) error {
	return merr.NewCommonError(modelNotLoadedErrorID, merr.ErrorKindValidate, message, nil)
}

// newSavePathInvalidError は保存先パス不正エラーを生成する。
func newSavePathInvalidError(message string) error {
	return merr.NewCommonError(savePathInvalidErrorID, merr.ErrorKindValidate, message, nil)
}

// logLoadFailed は読み込み失敗ログを出力する。
func logLoadFailed(logger logging.ILogger, translator i18n.II18n, err error) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	logErrorTitle(logger, translate(translator, "読み込み失敗"), err)
}

// logSaveFailed は保存失敗ログを出力する。
func logSaveFailed(logger logging.ILogger, translator i18n.II18n, err error) {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	logErrorTitle(logger, translate(translator, "保存失敗"), err)
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
	logger.Error("%s: %s", title, err.Error())
}

// isPmxConvertiblePath はPMX保存対象のパスか判定する。
func isPmxConvertiblePath(path string) bool {
	if path == "" {
		return false
	}
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, ".x") || strings.EqualFold(ext, ".pmd")
}

// buildPmxOutputPath はX/PMDファイルと同じ場所にPMX出力パスを生成する。
func buildPmxOutputPath(path string) string {
	dir, name, _ := mfile.SplitPath(path)
	if dir == "" || name == "" {
		return ""
	}
	return mfile.CreateOutputPath(filepath.Join(dir, name+".pmx"), "")
}

// translate は翻訳済み文言を返す。
func translate(translator i18n.II18n, key string) string {
	if translator == nil || !translator.IsReady() {
		return "●●" + key + "●●"
	}
	return translator.T(key)
}
