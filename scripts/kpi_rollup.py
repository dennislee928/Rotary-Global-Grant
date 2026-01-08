#!/usr/bin/env python3
"""KPI rollup template.

Usage:
  python3 scripts/kpi_rollup.py docs/03_kpi_and_me/kpi_template.csv

This script is intentionally minimal — it demonstrates how to parse KPI definitions.
Extend it to read workshop/incident exports and compute monthly rollups.
"""

import csv
import sys

def main(path: str) -> int:
    with open(path, newline="", encoding="utf-8") as f:
        rows = list(csv.DictReader(f))
    print(f"Loaded {len(rows)} KPI definitions")
    for r in rows:
        print(f"- {r['category']}/{r['metric']}: target={r['target_12m']}")
    return 0

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: kpi_rollup.py <kpi_template.csv>")
        raise SystemExit(2)
    raise SystemExit(main(sys.argv[1]))
