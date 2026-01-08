#!/bin/bash
# Generate TypeScript types from OpenAPI spec
#
# Prerequisites:
#   npm install -g openapi-typescript
#
# Usage:
#   ./scripts/generate_types.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

OPENAPI_PATH="$PROJECT_ROOT/packages/openapi/openapi.yaml"
OUTPUT_PATH="$PROJECT_ROOT/apps/web/lib/generated-types.ts"

echo "Generating TypeScript types from OpenAPI spec..."

if ! command -v npx &> /dev/null; then
    echo "Error: npx not found. Please install Node.js"
    exit 1
fi

npx openapi-typescript "$OPENAPI_PATH" -o "$OUTPUT_PATH"

echo "Generated types at: $OUTPUT_PATH"
