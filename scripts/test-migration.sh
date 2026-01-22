#!/bin/bash
set -e

echo "=== Testing Database Migration with 1Password CLI ==="
echo ""

# Check if op CLI is available
if ! command -v op &> /dev/null; then
    echo "Error: 1Password CLI (op) is not installed"
    echo "Install from: https://developer.1password.com/docs/cli/get-started/"
    exit 1
fi

echo "✓ 1Password CLI found"
echo ""

# Test reading credentials
echo "Testing credential access..."
if op read "op://pedro/POSTGRES_URL/credential" &> /dev/null; then
    echo "✓ Can access POSTGRES_URL from 1Password"
else
    echo "✗ Cannot access POSTGRES_URL from 1Password"
    echo "  Make sure you're signed in: op signin"
    exit 1
fi

echo ""
echo "Testing database connection..."
echo ""

# Test database connection using psql
POSTGRES_URL=$(op read "op://pedro/POSTGRES_URL/credential")

if psql "$POSTGRES_URL" -c "SELECT version FROM goose_db_version ORDER BY id DESC LIMIT 1;" 2>/dev/null; then
    echo ""
    echo "✓ Database connection successful"
    echo "✓ Migrations tracking table exists"
else
    echo ""
    echo "Note: Database may not be initialized yet"
    echo "Migrations will run automatically on first startup"
fi

echo ""
echo "=== Test Complete ==="
echo ""
echo "To run the application (migrations run automatically):"
echo "  op run --env-file=.env.op -- go run main.go"
echo ""
echo "To manually check migration status (optional):"
echo "  goose -dir db/migrations postgres \"\$POSTGRES_URL\" status"
