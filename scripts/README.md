# Scripts

Utility scripts for the Hive x Anti-Fraud project.

## Available Scripts

### kpi_rollup.py

KPI rollup and reporting tool. Loads KPI definitions from CSV and generates reports.

```bash
# Basic usage - print summary
python3 scripts/kpi_rollup.py docs/03_kpi_and_me/kpi_template.csv

# Export to JSON
python3 scripts/kpi_rollup.py --json output.json

# Generate markdown report
python3 scripts/kpi_rollup.py --markdown > kpi_report.md
```

### export_reports.py

Export reports from the API to CSV format for sponsor reporting.

```bash
# Export from local API
python3 scripts/export_reports.py --output reports.csv

# Export from production with auth
python3 scripts/export_reports.py \
  --api https://api.hive.example.invalid \
  --token $API_TOKEN \
  --output reports.csv \
  --summary
```

### generate_types.sh

Generate TypeScript types from the OpenAPI specification.

```bash
# Requires: npm install -g openapi-typescript
./scripts/generate_types.sh
```

## Development

These scripts are designed to be simple and dependency-light. They use only
Python standard library (for Python scripts) to avoid additional setup.

For production use, consider:
- Adding proper error handling
- Implementing retry logic for API calls
- Adding logging
- Using a proper HTTP client library
