#!/bin/bash
set -e

# リリーススクリプト
# 使い方: ./release.sh <major|minor|patch>

if [ $# -ne 1 ]; then
    echo "Usage: $0 <major|minor|patch>"
    echo "Example: $0 patch"
    exit 1
fi

# 現在のバージョンを取得
CURRENT_VERSION=$(grep 'const Version' version.go | sed 's/.*"\(.*\)".*/\1/')
echo "Current version: $CURRENT_VERSION"

# セマンティックバージョニングで次のバージョンを計算
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]}"
PATCH="${VERSION_PARTS[2]}"

case "$1" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Invalid version type: $1"
        echo "Use: major, minor, or patch"
        exit 1
        ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
echo "New version: $NEW_VERSION"

# リリースを実行
make release VERSION="$NEW_VERSION"