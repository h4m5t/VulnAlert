#!/bin/bash

OUTPUT_DIR="dist"
mkdir -p $OUTPUT_DIR

export CGO_ENABLED=0

declare -a TARGETS=(
    "windows/amd64"
    "linux/amd64"
    "darwin/arm64"
)

for TARGET in "${TARGETS[@]}"
do
    GOOS=$(echo $TARGET | cut -d '/' -f 1)
    GOARCH=$(echo $TARGET | cut -d '/' -f 2)
    
    OUTPUT_NAME="watchvuln-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME+=".exe"
    fi
    
    echo "编译目标：$GOOS/$GOARCH -> $OUTPUT_NAME"
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT_DIR/$OUTPUT_NAME main.go
    
    if [ $? -eq 0 ]; then
        echo "成功编译：$OUTPUT_NAME"
    else
        echo "编译失败：$OUTPUT_NAME"
    fi
done

echo "所有编译完成，输出位于 '$OUTPUT_DIR' 文件夹中。"