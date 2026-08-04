#!/usr/bin/env bash
# Syncs local wiki/ folder to GitHub Wiki remote repository (https://github.com/arkalon76/halptask.wiki.git)

set -e

WIKI_DIR="wiki"
TEMP_DIR=$(mktemp -d)
REMOTE_REPO="git@github.com:arkalon76/halptask.wiki.git"

echo "🚀 Syncing HalpTask documentation in '${WIKI_DIR}' to GitHub Wiki (${REMOTE_REPO})..."

if [ ! -d "$WIKI_DIR" ]; then
    echo "❌ Error: '${WIKI_DIR}' directory not found!"
    exit 1
fi

# Clone or init temp wiki repo
git clone "$REMOTE_REPO" "$TEMP_DIR" 2>/dev/null || (
    cd "$TEMP_DIR"
    git init
    git remote add origin "$REMOTE_REPO"
)

# Copy wiki markdown and images
cp -r "$WIKI_DIR"/* "$TEMP_DIR"/

cd "$TEMP_DIR"
git add .
if git diff-index --quiet HEAD -- 2>/dev/null; then
    echo "✨ Wiki is already up to date."
else
    git commit -m "docs(wiki): sync wiki documentation from main repository"
    git branch -M main 2>/dev/null || true
    git push origin main || git push origin master
    echo "✅ GitHub Wiki sync completed successfully!"
fi

rm -rf "$TEMP_DIR"
