// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mu_model_viewer/pkg/usecase/port/moutput"
)

// ModelLoadResult はモデル読み込み結果を表す。
type ModelLoadResult = usecase.ModelLoadResult

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult = usecase.MotionLoadResult

// TextureValidationResult はテクスチャ検証結果を表す。
type TextureValidationResult = usecase.TextureValidationResult

// PmxSaveResult はPMX保存結果を表す。
type PmxSaveResult = usecase.PmxSaveResult

// SaveOptions は保存時のオプションを表す。
type SaveOptions = moutput.SaveOptions
