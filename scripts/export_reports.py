#!/usr/bin/env python3
"""Export reports from the API to CSV format.

Usage:
  python3 scripts/export_reports.py --output reports.csv
  python3 scripts/export_reports.py --api http://localhost:8080 --token $API_TOKEN

This script fetches reports from the API and exports them to CSV format
for sponsor reporting and analysis.
"""

import argparse
import csv
import json
import sys
from datetime import datetime
from typing import Optional
from urllib.request import Request, urlopen
from urllib.error import URLError, HTTPError


def fetch_reports(api_base: str, token: Optional[str] = None, page_size: int = 100) -> list[dict]:
    """Fetch all reports from the API with pagination."""
    reports = []
    page = 1
    
    while True:
        url = f"{api_base}/v1/reports?page={page}&pageSize={page_size}"
        
        headers = {"Content-Type": "application/json"}
        if token:
            headers["Authorization"] = f"Bearer {token}"
        
        try:
            req = Request(url, headers=headers)
            with urlopen(req) as response:
                data = json.loads(response.read().decode())
        except HTTPError as e:
            if e.code == 401:
                print("Error: Authentication required. Use --token option.", file=sys.stderr)
            else:
                print(f"Error fetching reports: {e}", file=sys.stderr)
            return reports
        except URLError as e:
            print(f"Error connecting to API: {e}", file=sys.stderr)
            return reports
        
        reports.extend(data.get("data", []))
        
        pagination = data.get("pagination", {})
        if page >= pagination.get("totalPages", 1):
            break
        page += 1
    
    return reports


def export_to_csv(reports: list[dict], output_path: str) -> None:
    """Export reports to CSV file."""
    if not reports:
        print("No reports to export.", file=sys.stderr)
        return
    
    fieldnames = [
        "id",
        "category",
        "severitySuggested",
        "areaHint",
        "timeWindow",
        "description",
        "status",
        "createdAt",
        "updatedAt",
    ]
    
    with open(output_path, "w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames, extrasaction="ignore")
        writer.writeheader()
        
        for report in reports:
            # Clean up the report for CSV export
            row = {k: report.get(k, "") for k in fieldnames}
            # Truncate long descriptions
            if len(row.get("description", "")) > 500:
                row["description"] = row["description"][:497] + "..."
            writer.writerow(row)
    
    print(f"Exported {len(reports)} reports to {output_path}")


def generate_summary(reports: list[dict]) -> dict:
    """Generate summary statistics from reports."""
    summary = {
        "total": len(reports),
        "by_status": {},
        "by_category": {},
        "date_range": {
            "earliest": None,
            "latest": None,
        },
    }
    
    for report in reports:
        # Count by status
        status = report.get("status", "unknown")
        summary["by_status"][status] = summary["by_status"].get(status, 0) + 1
        
        # Count by category
        category = report.get("category", "unknown")
        summary["by_category"][category] = summary["by_category"].get(category, 0) + 1
        
        # Track date range
        created_at = report.get("createdAt", "")
        if created_at:
            if not summary["date_range"]["earliest"] or created_at < summary["date_range"]["earliest"]:
                summary["date_range"]["earliest"] = created_at
            if not summary["date_range"]["latest"] or created_at > summary["date_range"]["latest"]:
                summary["date_range"]["latest"] = created_at
    
    return summary


def main() -> int:
    parser = argparse.ArgumentParser(description="Export reports from API to CSV")
    parser.add_argument("--api", default="http://localhost:8080", help="API base URL")
    parser.add_argument("--token", help="API authentication token")
    parser.add_argument("--output", "-o", default="reports.csv", help="Output CSV file")
    parser.add_argument("--summary", action="store_true", help="Print summary statistics")
    
    args = parser.parse_args()
    
    print(f"Fetching reports from {args.api}...")
    reports = fetch_reports(args.api, args.token)
    
    if not reports:
        print("No reports found or unable to fetch reports.")
        return 1
    
    if args.summary:
        summary = generate_summary(reports)
        print("\nSummary:")
        print(f"  Total reports: {summary['total']}")
        print(f"  Date range: {summary['date_range']['earliest']} to {summary['date_range']['latest']}")
        print("\n  By status:")
        for status, count in sorted(summary["by_status"].items()):
            print(f"    {status}: {count}")
        print("\n  By category:")
        for category, count in sorted(summary["by_category"].items()):
            print(f"    {category}: {count}")
    
    export_to_csv(reports, args.output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
