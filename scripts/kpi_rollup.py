#!/usr/bin/env python3
"""KPI rollup and reporting tool.

Usage:
  python3 scripts/kpi_rollup.py docs/03_kpi_and_me/kpi_template.csv
  python3 scripts/kpi_rollup.py --json output.json
  python3 scripts/kpi_rollup.py --markdown

This script parses KPI definitions and computes monthly rollups from
workshop/incident exports. Extend it to connect to the API for live data.
"""

import argparse
import csv
import json
import sys
from dataclasses import dataclass
from typing import Optional
from datetime import datetime


@dataclass
class KPIDefinition:
    category: str
    metric: str
    definition: str
    source: str
    frequency: str
    baseline: str
    target_12m: str
    owner: str


@dataclass
class KPIResult:
    category: str
    metric: str
    current: float
    target: float
    met: bool
    percentage: float


def load_kpi_definitions(path: str) -> list[KPIDefinition]:
    """Load KPI definitions from CSV file."""
    definitions = []
    with open(path, newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        for row in reader:
            definitions.append(KPIDefinition(
                category=row.get("category", ""),
                metric=row.get("metric", ""),
                definition=row.get("definition", ""),
                source=row.get("source", ""),
                frequency=row.get("frequency", ""),
                baseline=row.get("baseline", "0"),
                target_12m=row.get("target_12m", ""),
                owner=row.get("owner", ""),
            ))
    return definitions


def parse_target(target_str: str) -> tuple[float, bool]:
    """Parse target string like '>=300' or '<=5%'.
    Returns (value, is_inverse) where is_inverse means lower is better.
    """
    target_str = target_str.strip()
    is_inverse = False
    
    if target_str.startswith("<="):
        is_inverse = True
        target_str = target_str[2:]
    elif target_str.startswith(">="):
        target_str = target_str[2:]
    elif target_str.startswith("<"):
        is_inverse = True
        target_str = target_str[1:]
    elif target_str.startswith(">"):
        target_str = target_str[1:]
    
    # Remove % sign and other suffixes
    target_str = target_str.replace("%", "").replace("min*", "").replace("min", "").strip()
    
    try:
        value = float(target_str)
    except ValueError:
        value = 0.0
    
    return value, is_inverse


def compute_kpi_status(current: float, target: float, is_inverse: bool) -> tuple[bool, float]:
    """Compute if target is met and percentage progress."""
    if is_inverse:
        met = current <= target
        if target > 0:
            # For inverse metrics, 100% means at or below target
            percentage = max(0, min(100, ((target - current) / target) * 100 + 100))
        else:
            percentage = 100 if current == 0 else 0
    else:
        met = current >= target
        if target > 0:
            percentage = min(100, (current / target) * 100)
        else:
            percentage = 100 if current > 0 else 0
    
    return met, percentage


def generate_markdown_report(definitions: list[KPIDefinition], results: dict[str, float]) -> str:
    """Generate a markdown report from KPI definitions and results."""
    lines = [
        "# KPI Report",
        f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M')}",
        "",
        "## Summary",
        "",
        "| Category | Metric | Current | Target | Status |",
        "|----------|--------|---------|--------|--------|",
    ]
    
    for defn in definitions:
        key = f"{defn.category}/{defn.metric}"
        current = results.get(key, 0)
        target, is_inverse = parse_target(defn.target_12m)
        met, percentage = compute_kpi_status(current, target, is_inverse)
        
        status = "✅" if met else "🔄"
        target_display = defn.target_12m
        
        lines.append(
            f"| {defn.category} | {defn.metric} | {current:.1f} | {target_display} | {status} {percentage:.0f}% |"
        )
    
    lines.extend([
        "",
        "## Definitions",
        "",
    ])
    
    for defn in definitions:
        lines.extend([
            f"### {defn.category} / {defn.metric}",
            "",
            f"- **Definition**: {defn.definition}",
            f"- **Source**: {defn.source}",
            f"- **Frequency**: {defn.frequency}",
            f"- **Target**: {defn.target_12m}",
            f"- **Owner**: {defn.owner}",
            "",
        ])
    
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser(description="KPI rollup and reporting tool")
    parser.add_argument("kpi_file", nargs="?", default="docs/03_kpi_and_me/kpi_template.csv",
                        help="Path to KPI definitions CSV")
    parser.add_argument("--json", metavar="FILE", help="Output results as JSON")
    parser.add_argument("--markdown", action="store_true", help="Output markdown report")
    parser.add_argument("--api", metavar="URL", help="Fetch live data from API")
    
    args = parser.parse_args()
    
    # Load definitions
    try:
        definitions = load_kpi_definitions(args.kpi_file)
        print(f"Loaded {len(definitions)} KPI definitions from {args.kpi_file}")
    except FileNotFoundError:
        print(f"Error: File not found: {args.kpi_file}", file=sys.stderr)
        return 1
    except Exception as e:
        print(f"Error loading KPI file: {e}", file=sys.stderr)
        return 1
    
    # Demo results (in production, fetch from API)
    results = {
        "education/participants_trained": 342,
        "education/pre_post_improvement": 25.5,
        "education/workshops_count": 12,
        "governance/trained_triagers": 18,
        "system/median_report_to_triage": 22.5,
        "system/verified_ratio": 65.2,
        "system/abuse_rate": 3.2,
        "alerts/publish_latency": 12.0,
        "accessibility/a11y_language_coverage": 1,
        "adoption/partner_orgs": 4,
        "adoption/external_adoption": 2,
    }
    
    if args.api:
        print(f"Note: API fetching not implemented. Using demo data.")
    
    # Output
    if args.json:
        output = {
            "generated_at": datetime.now().isoformat(),
            "definitions": [vars(d) for d in definitions],
            "results": results,
        }
        with open(args.json, "w") as f:
            json.dump(output, f, indent=2)
        print(f"Wrote JSON output to {args.json}")
    elif args.markdown:
        report = generate_markdown_report(definitions, results)
        print(report)
    else:
        # Default: print summary
        print("\nKPI Summary:")
        print("-" * 60)
        for defn in definitions:
            key = f"{defn.category}/{defn.metric}"
            current = results.get(key, 0)
            target, is_inverse = parse_target(defn.target_12m)
            met, percentage = compute_kpi_status(current, target, is_inverse)
            
            status = "✓" if met else "○"
            print(f"  {status} {defn.category}/{defn.metric}: {current:.1f} (target: {defn.target_12m})")
    
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
