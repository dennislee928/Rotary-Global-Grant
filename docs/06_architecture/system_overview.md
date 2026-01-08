# System Overview (Suggested Architecture)

## Goals
- Safe intake of community reports
- Fast triage + verification workflow
- Tiered alerts with consistent messaging (CAP-ready)
- Education hub (anti-fraud + crisis communication)
- Strong auditability and abuse prevention

## Suggested components
1) Web UI (Next.js)
   - reporting form
   - public dashboard (low-risk info only)
   - training modules + quizzes
2) API (Go/Gin)
   - auth (optional for reporters)
   - report intake + evidence handling
   - triage decisions + audit log
   - CAP alert composer and publishing workflow
3) Data layer
   - Postgres: transactional data (reports, triage, alerts, training outcomes)
   - Redis: queues, rate limiters, ephemeral tokens
4) Observability
   - structured logs
   - metrics endpoints for KPI rollups
5) Integrations (optional)
   - email/SMS gateway
   - LINE/WhatsApp bots
   - export to CSV for sponsor reporting

See `docs/06_architecture/data_model.md` and `packages/openapi/openapi.yaml`.
