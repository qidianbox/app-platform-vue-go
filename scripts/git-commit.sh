#!/bin/bash

# 自动提交脚本
# 用法: ./scripts/git-commit.sh "提交信息"

set -e

# 检查是否提供了提交信息
if [ -z "$1" ]; then
    echo "错误: 请提供提交信息"
    echo "用法: ./scripts/git-commit.sh \"提交信息\""
    exit 1
fi

COMMIT_MESSAGE="$1"
PROJECT_DIR="/home/ubuntu/app-platform-vue-go"

cd "$PROJECT_DIR"

# 配置 Git 用户信息（如果未配置）
git config user.name "Manus AI" 2>/dev/null || true
git config user.email "ai@manus.im" 2>/dev/null || true

# 显示当前状态
echo "=== Git 状态 ==="
git status --short

# 添加所有更改
echo ""
echo "=== 添加更改 ==="
git add -A

# 显示将要提交的文件
echo ""
echo "=== 将要提交的文件 ==="
git diff --cached --name-status

# 提交
echo ""
echo "=== 提交更改 ==="
git commit -m "$COMMIT_MESSAGE" || {
    echo "提示: 没有需要提交的更改"
    exit 0
}

# 推送到远程仓库
echo ""
echo "=== 推送到 GitHub ==="
git push origin master

echo ""
echo "✅ 成功提交并推送到 GitHub!"
echo "提交信息: $COMMIT_MESSAGE"
