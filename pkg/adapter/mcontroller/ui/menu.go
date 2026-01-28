//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"

	"github.com/miu200521358/mu_model_viewer/pkg/adapter/mpresenter/messages"
)

// NewMenuItems は mu_model_viewer のメニュー項目を生成する。
func NewMenuItems(translator i18n.II18n, logger logging.ILogger) []declarative.MenuItem {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	t := func(key string) string {
		if translator == nil || !translator.IsReady() {
			return "●●" + key + "●●"
		}
		return translator.T(key)
	}
	return []declarative.MenuItem{
		declarative.Action{
			Text: t(messages.HelpUsage),
			OnTriggered: func() {
				if lineLogger, ok := logger.(interface {
					InfoLine(msg string, params ...any)
				}); ok {
					lineLogger.InfoLine(messages.HelpUsage)
					return
				}
				logger.Info(messages.HelpUsage)
			},
		},
		declarative.Action{
			Text: t(messages.HelpMaterialView),
			OnTriggered: func() {
				if lineLogger, ok := logger.(interface {
					InfoLine(msg string, params ...any)
				}); ok {
					lineLogger.InfoLine(messages.HelpMaterialView)
					return
				}
				logger.Info(messages.HelpMaterialView)
			},
		},
	}
}
