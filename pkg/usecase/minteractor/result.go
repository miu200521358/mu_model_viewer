// 指示: miu200521358
package minteractor

import (
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase/port/moutput"
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
type SaveOptions = moutput.SaveOptions
