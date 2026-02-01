// 指示: miu200521358
package minteractor

import (
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase/port/moutput"
)

// LoadModel はモデルを読み込み、テクスチャ検証結果を付与して返す。
func (uc *ModelViewerUsecase) LoadModel(rep moutput.IFileReader, path string) (*ModelLoadResult, error) {
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
func (uc *ModelViewerUsecase) LoadMotion(rep moutput.IFileReader, path string) (*MotionLoadResult, error) {
	repo := rep
	if repo == nil {
		repo = uc.motionReader
	}
	return commonusecase.LoadMotionWithMeta(repo, path)
}

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func (uc *ModelViewerUsecase) CanLoadPath(rep moutput.IFileReader, path string) bool {
	repo := rep
	if repo == nil {
		repo = uc.modelReader
	}
	return commonusecase.CanLoadPath(repo, path)
}
