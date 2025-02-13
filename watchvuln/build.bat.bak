@echo off
chcp 65001
setlocal enabledelayedexpansion

set OUTPUT_DIR=dist
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

set CGO_ENABLED=0

:: 定义目标平台
set TARGETS=windows/amd64 linux/amd64 darwin/arm64

for %%t in (%TARGETS%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%t") do (
        set GOOS=%%a
        set GOARCH=%%b
        
        set OUTPUT_NAME=watchvuln-!GOOS!-!GOARCH!
        if "!GOOS!"=="windows" (
            set OUTPUT_NAME=!OUTPUT_NAME!.exe
        )
        
        echo [编译目标]：!GOOS!/!GOARCH! -^> !OUTPUT_NAME!
        
        set GOOS=!GOOS!
        set GOARCH=!GOARCH!
        go build -o %OUTPUT_DIR%/!OUTPUT_NAME! main.go
        
        if !errorlevel! equ 0 (
            echo [成功]：!OUTPUT_NAME!
        ) else (
            echo [失败]：!OUTPUT_NAME!
        )
    )
)

echo [完成] 所有编译完成，输出位于 '%OUTPUT_DIR%' 文件夹中。
pause