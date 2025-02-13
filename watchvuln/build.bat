@echo off
chcp 65001 > nul
setlocal enabledelayedexpansion

set "OUTPUT_DIR=dist"
if not exist "%OUTPUT_DIR%" (
    mkdir "%OUTPUT_DIR%" || (
        echo 无法创建目录 "%OUTPUT_DIR%"
        pause
        exit /b 1
    )
)

set "CGO_ENABLED=0"

:: 定义目标平台（使用逗号分隔，避免空格问题）
set "TARGETS=windows/amd64,linux/amd64,darwin/arm64"

for %%t in (%TARGETS%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%t") do (
        set "GOOS=%%a"
        set "GOARCH=%%b"
        
        set "OUTPUT_NAME=watchvuln-!GOOS!-!GOARCH!"
        if "!GOOS!"=="windows" (
            set "OUTPUT_NAME=!OUTPUT_NAME!.exe"
        )
        
        echo [编译目标]: !GOOS!/!GOARCH! ^> "!OUTPUT_NAME!"
        
        :: 显式设置环境变量
        set "GOOS=!GOOS!"
        set "GOARCH=!GOARCH!"
        
        :: 编译命令（路径用双引号包裹）
        echo 正在编译: go build -o "%OUTPUT_DIR%\!OUTPUT_NAME!" main.go
        go build -o "%OUTPUT_DIR%\!OUTPUT_NAME!" main.go
        
        if !errorlevel! neq 0 (
            echo [错误]: 编译 "!OUTPUT_NAME!" 失败，错误码=!errorlevel!
            pause
            exit /b !errorlevel!
        )
    )
)

echo [完成] 所有编译完成，输出位于 "%OUTPUT_DIR%" 文件夹中。
pause