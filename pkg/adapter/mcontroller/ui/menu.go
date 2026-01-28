//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"

	"github.com/miu200521358/mu_model_viewer/pkg/adapter/mpresenter/messages"
)

// NewMenuItems は mu_model_viewer のメニュー項目を生成する。
func NewMenuItems(translator i18n.II18n, logger logging.ILogger) []declarative.MenuItem {
	return controller.BuildMenuItemsWithMessages(translator, logger, []controller.MenuMessageItem{
		{TitleKey: messages.HelpUsageTitle, MessageKey: messages.HelpUsage},
		{TitleKey: controller.MenuSeparatorKey},
		{TitleKey: messages.HelpMaterialViewTitle, MessageKey: messages.HelpMaterialView},
	})
}
