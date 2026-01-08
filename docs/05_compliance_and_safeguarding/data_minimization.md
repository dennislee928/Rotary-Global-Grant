# Data Minimization & Privacy by Design (Draft)

## Principles
- Collect the minimum data needed for triage
- Avoid precise home addresses, national IDs, and sensitive identifiers
- Prefer approximate area and time windows
- Separate “contactable” data from incident data (pseudonymous keys)

## Retention
- Define retention windows per data class:
  - raw reports: X days
  - audit logs: Y months
  - aggregated metrics: Z years (no personal identifiers)

## Access control
- Role-based access (triage vs admin vs auditor)
- Least privilege
- Immutable audit log for decisions
