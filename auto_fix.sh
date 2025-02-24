#!/bin/bash

# Set the code directory
CODE_DIR="/github/workspace"

# Apply clang-tidy fixes
echo "Running clang-tidy with autofix..."
find "$CODE_DIR" -type f \( -name "*.cpp" -o -name "*.h" \) -exec clang-tidy -fix {} -- -std=c++17 \;

# Verify no new issues
echo "Re-running clang-tidy to confirm fixes..."
find "$CODE_DIR" -type f \( -name "*.cpp" -o -name "*.h" \) -exec clang-tidy {} -- -std=c++17 \;

echo "Clang-tidy fixes applied successfully!"

