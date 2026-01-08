# The Hive x Anti-Fraud — Rotary Global Grant Scholarship Toolkit (Repo Template)

This repository is a **ready-to-adapt template** for proposing and executing a Rotary-aligned project that combines:

- **The Hive** (community safety reporting + verification + tiered alerts, CAP-ready)
- **Anti-Fraud Platform** (anti-scam education, playbooks, workshops, and evaluation)

It includes:
- Rotary process guidance (how to “apply” in practice)
- Community Needs Assessment (Taipei / Glasgow versions)
- KPI measurement framework and M&E SOP
- Pilot playbooks (tabletop + drill) that you can run immediately
- A practical, buildable **reference code architecture** (Go + Postgres + Redis + Next.js)

> Status: scaffold (draft). Last updated: 2026-01-08

---

## Quick start

### 1) Use the docs (proposal + sponsor conversations)
Start here:

- `docs/00_project_brief_one_pager.md`
- `docs/01_rotary_process.md`
- `docs/02_community_needs_assessment/`
- `docs/03_kpi_and_me/`
- `docs/04_pilot_playbooks/`
- `docs/08_timeline/roadmap_12_months.md`

### 2) Run the reference stack locally (optional)
You can run a minimal reference stack with Docker Compose:

```bash
cd infra
docker compose up --build
```

- Web: http://localhost:3000
- API: http://localhost:8080/healthz

> The API is designed to run with Postgres/Redis but can also run in “dev-memory” mode.

---

## Intended Rotary alignment

This template is written to support proposals aligned to:
- **Promoting peace** (community safety, de-escalation, panic reduction, trusted communications)
- **Supporting education** (anti-fraud + digital safety literacy, train-the-trainer, community workshops)

See `docs/00_project_brief_one_pager.md`.

---

## Repository layout

```
.
├─ apps/
│  ├─ api/                 # Go (Gin) reference API
│  └─ web/                 # Next.js reference web UI
├─ docs/                   # Proposal + research + playbooks
├─ infra/                  # docker-compose + env templates
├─ packages/
│  └─ openapi/             # OpenAPI spec & schemas
└─ scripts/                # helper scripts (KPI export templates)
```

---

## License

Choose a license that matches your intent (MIT/Apache-2.0 recommended for open-source).  
A placeholder is included at `LICENSE` — update as needed.
