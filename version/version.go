package version

import "github.com/Droi-SDK/droi-checker/logger"

const Version = "0.1.1"

func PrintCurrentVersion() {
	logger.Info("当前命令行工具版本：", Version)
}
