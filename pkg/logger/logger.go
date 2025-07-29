package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

var (
	infoLogger  *slog.Logger
	writeLogger *slog.Logger
	info        *os.File
	write       *os.File
)

// 用来获取调用函数信息方便打印
func caller() slog.Attr {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return slog.String("caller", "unknown")
	}
	return slog.String("caller", filepath.Base(file)+":"+string(rune(line)))
}
func InitLogger() {
	//创建log文件夹
	os.Mkdir("log", 0755)
	//创建info日志文件
	infoFile, _ := os.OpenFile("log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	infoHandler := slog.NewJSONHandler(infoFile, nil)
	infoLogger = slog.New(infoHandler)
	//创建write日志文件
	writeFile, _ := os.OpenFile("log/root.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writeHandler := slog.NewTextHandler(writeFile, nil)
	writeLogger = slog.New(writeHandler)
}

func Info(msg string, args ...any) {
	infoLogger.Info(msg, append(args, caller())...)
}

func Debug(msg string, args ...any) {
	infoLogger.Debug(msg, append(args, caller())...)
}

func Error(msg string, args ...any) {
	infoLogger.Error(msg, append(args, caller())...)
}
func Write(msg string, args ...any) {
	writeLogger.Info(msg, append(args, caller())...)
}

func Close() {
	if info != nil {
		info.Close()
	}
	if write != nil {
		write.Close()
	}
}
