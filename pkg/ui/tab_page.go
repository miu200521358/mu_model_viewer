package ui

import (
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/mfile"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewTabPage(mWidgets *controller.MWidgets) declarative.TabPage {
	var fileTab *walk.TabPage

	player := widget.NewMotionPlayer()

	materialTableView := widget.NewMaterialTableView(
		mi18n.T("材質ビュー説明"),
		func(cw *controller.ControlWindow, indexes []int) {
			cw.StoreSelectedMaterialIndexes(0, 0, indexes)
		},
	)

	pmxLoadPicker := widget.NewPmxXLoadFilePicker(
		"pmx",
		mi18n.T("モデルファイル"),
		mi18n.T("モデルファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if path == "" {
				cw.StoreModel(0, 0, nil)
				return
			}

			if data, err := rep.Load(path); err == nil {
				model := data.(*pmx.PmxModel)
				cw.StoreModel(0, 0, model)

				model.Textures.ForEach(func(index int, texture *pmx.Texture) bool {
					// モデルパス + テクスチャ相対パス
					texPath := filepath.Join(filepath.Dir(model.Path()), texture.Name())

					// テクスチャの有効判定
					texture.SetValid(true)

					valid, err := mfile.ExistsFile(texPath)
					if !valid || err != nil {
						texture.SetValid(false)
					}

					// 画像を読み込み
					if _, err := mfile.LoadImage(texPath); err != nil {
						texture.SetValid(false)
					}

					return true
				})

				// モデルの読み込みが成功したら材質テーブル更新
				materialTableView.MaterialModel.ResetRows(model)

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
			if path == "" {
				cw.StoreMotion(0, 0, nil)
				return
			}

			if data, err := rep.Load(path); err == nil {
				motion := data.(*vmd.VmdMotion)
				player.Reset(motion.MaxFrame())
				cw.StoreMotion(0, 0, motion)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
		},
	)

	// 全選択ボタン
	allBtn := widget.NewMPushButton()
	allBtn.SetLabel(mi18n.T("全"))
	allBtn.SetTooltip(mi18n.T("全ボタン説明"))
	allBtn.SetMaxSize(declarative.Size{Width: 50})
	allBtn.SetOnClicked(func(cw *controller.ControlWindow) {
		for _, record := range materialTableView.MaterialModel.Records {
			record.Checked = true
		}
		materialTableView.MaterialModel.PublishRowsChanged(0, materialTableView.MaterialModel.RowCount())
		cw.StoreSelectedMaterialIndexes(0, 0, materialTableView.MaterialModel.CheckedIndexes())
	})

	// 選択反転ボタン
	revertBtn := widget.NewMPushButton()
	revertBtn.SetLabel(mi18n.T("反"))
	revertBtn.SetTooltip(mi18n.T("反ボタン説明"))
	revertBtn.SetMaxSize(declarative.Size{Width: 50})
	revertBtn.SetOnClicked(func(cw *controller.ControlWindow) {
		for _, record := range materialTableView.MaterialModel.Records {
			record.Checked = !record.Checked
		}
		materialTableView.MaterialModel.PublishRowsChanged(0, materialTableView.MaterialModel.RowCount())
		cw.StoreSelectedMaterialIndexes(0, 0, materialTableView.MaterialModel.CheckedIndexes())
	})

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoadPicker, vmdLoadPicker, materialTableView, allBtn, revertBtn)
	mWidgets.SetOnLoaded(func() {
		// 読み込みが完了したら、モデルのパスを設定
		if path, err := usecase.LoadModelPath(); err == nil {
			pmxLoadPicker.SetPath(path)
		}
	})

	return declarative.TabPage{
		Title:    mi18n.T("ファイル"),
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
							declarative.TextLabel{
								Text: mi18n.T("材質ビュー"),
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									mlog.I(mi18n.T("材質ビュー説明"))
								},
							},
							declarative.HSpacer{},
							allBtn.Widgets(),
							revertBtn.Widgets(),
						},
					},
					materialTableView.Widgets(),
					declarative.VSpacer{},
					player.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}
}
