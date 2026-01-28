// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
)

// SaveModelAsPmxRequest はPMX保存の入力を表す。
type SaveModelAsPmxRequest struct {
	ModelPath              string
	ModelData              *model.PmxModel
	MissingModelMessage    string
	InvalidSavePathMessage string
	SaveOptions            SaveOptions
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
