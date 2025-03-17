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

	materialListbox := widget.NewMaterialListbox(
		"", func(cw *controller.ControlWindow, indexes []int) {
			cw.StoreSelectedMaterialIndexes(0, 0, indexes)
		},
	)

	pmxLoadPicker := widget.NewPmxXLoadFilePicker(
		"pmx",
		mi18n.T("モデルファイル"),
		mi18n.T("モデルファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if data, err := rep.Load(path); err == nil {
				model := data.(*pmx.PmxModel)
				cw.StoreModel(0, 0, model)

				// モデルの読み込みが成功したら材質リスト更新
				materialListbox.SetMaterials(model.Materials)

				// フォーカスを当てる
				cw.SetFocus()
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	vmdLoadPicker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		mi18n.T("モーションファイル"),
		mi18n.T("モーションファイルを選択してください"),
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

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoadPicker, vmdLoadPicker, materialListbox)
	mWidgets.SetOnLoaded(func() {
		// 読み込みが完了したら、モデルのパスを設定
		if path, err := usecase.LoadModelPath(); err == nil {
			pmxLoadPicker.SetPath(path)
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
						pmxLoadPicker.Widgets(),
						vmdLoadPicker.Widgets(),
						declarative.VSeparator{},
						declarative.TextLabel{
							Text: mi18n.T("材質リスト"),
						},
						declarative.TextLabel{
							Text: mi18n.T("材質リスト説明"),
						},
						materialListbox.Widgets(),
						declarative.VSpacer{},
						player.Widgets(),
						declarative.VSpacer{},
					},
				},
			},
		},
	}
}
