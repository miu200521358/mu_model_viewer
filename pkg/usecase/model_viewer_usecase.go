// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	portio "github.com/miu200521358/mu_model_viewer/pkg/usecase/port/io"
)

// ModelLoadResult はモデル読み込み結果を表す。
type ModelLoadResult = commonusecase.ModelLoadResult

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult = commonusecase.MotionLoadResult

// TextureValidationResult はテクスチャ検証結果を表す。
type TextureValidationResult = commonusecase.TextureValidationResult

// PmxSaveResult はPMX保存結果を表す。
type PmxSaveResult = commonusecase.PmxSaveResult

// SaveOptions は保存時のオプションを表す。
type SaveOptions = portio.SaveOptions

// ModelViewerUsecaseDeps はモデルビューア用ユースケースの依存を表す。
type ModelViewerUsecaseDeps struct {
	ModelReader      portio.IFileReader
	MotionReader     portio.IFileReader
	ModelWriter      portio.IFileWriter
	TextureValidator portio.ITextureValidator
	PathService      portio.IPathService
}

// ModelViewerUsecase はモデルビューアの入出力処理をまとめたユースケースを表す。
type ModelViewerUsecase struct {
	modelReader      portio.IFileReader
	motionReader     portio.IFileReader
	modelWriter      portio.IFileWriter
	textureValidator portio.ITextureValidator
	pathService      portio.IPathService
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

// LoadModel はモデルを読み込み、テクスチャ検証結果を付与して返す。
func (uc *ModelViewerUsecase) LoadModel(rep portio.IFileReader, path string) (*ModelLoadResult, error) {
	repo := rep
	if repo == nil {
		repo = uc.modelReader
	}
	modelData, err := commonusecase.LoadModel(repo, path)
	if err != nil {
		return nil, err
	}
	result := &ModelLoadResult{Model: modelData}
	if modelData != nil && uc.textureValidator != nil {
		result.Validation = commonusecase.ValidateModelTextures(modelData, uc.textureValidator)
	}
	return result, nil
}

// LoadMotion はモーションを読み込み、最大フレーム情報を返す。
func (uc *ModelViewerUsecase) LoadMotion(rep portio.IFileReader, path string) (*MotionLoadResult, error) {
	repo := rep
	if repo == nil {
		repo = uc.motionReader
	}
	return commonusecase.LoadMotionWithMeta(repo, path)
}

// SaveModelAsPmxRequest はPMX保存の入力を表す。
type SaveModelAsPmxRequest struct {
	ModelPath              string
	ModelData              *model.PmxModel
	MissingModelMessage    string
	InvalidSavePathMessage string
	SaveOptions            portio.SaveOptions
}

// SaveModelAsPmx はX/PMDモデルをPMX形式で保存する。
func (uc *ModelViewerUsecase) SaveModelAsPmx(request SaveModelAsPmxRequest) (*PmxSaveResult, error) {
	return commonusecase.SaveModelAsPmx(commonusecase.PmxSaveRequest{
		ModelPath:              request.ModelPath,
		ModelData:              request.ModelData,
		Writer:                 uc.modelWriter,
		PathService:            uc.pathService,
		MissingModelMessage:    request.MissingModelMessage,
		InvalidSavePathMessage: request.InvalidSavePathMessage,
		SaveOptions:            request.SaveOptions,
	})
}

// IsPmxConvertiblePath はPMX保存対象のパスか判定する。
func (uc *ModelViewerUsecase) IsPmxConvertiblePath(path string) bool {
	return commonusecase.IsPmxConvertiblePath(path)
}

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func (uc *ModelViewerUsecase) CanLoadPath(rep portio.IFileReader, path string) bool {
	repo := rep
	if repo == nil {
		repo = uc.modelReader
	}
	return commonusecase.CanLoadPath(repo, path)
}
