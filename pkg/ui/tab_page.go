package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewTabPages(mWidgets *controller.MWidgets) []declarative.TabPage {
	var fileTab *walk.TabPage

	player := widget.NewMotionPlayer()

	pmxLoad11Picker := widget.NewPmxXLoadFilePicker(
		"pmx",
		"モデルファイル1-1",
		"モデルファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				cw.StoreModel(0, 0, data.(*pmx.PmxModel))
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	vmdLoader11Picker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		"モーションファイル1-1",
		"モーションファイルを選択してください",
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				motion := data.(*vmd.VmdMotion)
				player.Reset(motion.MaxFrame())
				cw.StoreMotion(0, 0, motion)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoad11Picker, vmdLoader11Picker)
	mWidgets.SetOnLoaded(func() {
		// 読み込みが完了したら、モデルのパスを設定
		if path, err := usecase.LoadModelPath(); err == nil {
			pmxLoad11Picker.SetPath(path)
		}
	})

	return []declarative.TabPage{
		{
			Title:    mi18n.T("ファイル"),
			AssignTo: &fileTab,
			Layout:   declarative.VBox{},
			Background: declarative.SystemColorBrush{
				Color: walk.SysColorInactiveCaption,
			},
			Children: []declarative.Widget{
				declarative.Composite{
					Layout: declarative.VBox{},
					Children: []declarative.Widget{
						declarative.TextLabel{
							Text: mi18n.T("ツール説明"),
						},
						declarative.VSeparator{},
						pmxLoad11Picker.Widgets(),
						vmdLoader11Picker.Widgets(),
						declarative.VSeparator{},
						player.Widgets(),
						declarative.VSpacer{},
					},
				},
			},
		},
	}
}
