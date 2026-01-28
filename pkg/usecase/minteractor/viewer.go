// 指示: miu200521358
package minteractor

import "github.com/miu200521358/mu_model_viewer/pkg/usecase/port/moutput"

// ModelViewerUsecaseDeps はモデルビューア用ユースケースの依存を表す。
type ModelViewerUsecaseDeps struct {
	ModelReader      moutput.IFileReader
	MotionReader     moutput.IFileReader
	ModelWriter      moutput.IFileWriter
	TextureValidator moutput.ITextureValidator
	PathService      moutput.IPathService
}

// ModelViewerUsecase はモデルビューアの入出力処理をまとめたユースケースを表す。
type ModelViewerUsecase struct {
	modelReader      moutput.IFileReader
	motionReader     moutput.IFileReader
	modelWriter      moutput.IFileWriter
	textureValidator moutput.ITextureValidator
	pathService      moutput.IPathService
}

// NewModelViewerUsecase はモデルビューア用ユースケースを生成する。
func NewModelViewerUsecase(deps ModelViewerUsecaseDeps) *ModelViewerUsecase {
	return &ModelViewerUsecase{
		modelReader:      deps.ModelReader,
		motionReader:     deps.MotionReader,
		modelWriter:      deps.ModelWriter,
		textureValidator: deps.TextureValidator,
		pathService:      deps.PathService,
	}
}
