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
___
Rotary Global Grant Scholarship 申請/送件線（用「T = 預計出發日」倒推）

12 個月實踐/交付線（The Hive x Anti-Fraud 的產品 + 教育 + pilot）

註：多個 Rotary 地區會有「地區內部截止日」早於基金會審查期；例如有地區明確公告需在秋季先交給社團/地區委員會審。
另：建議至少在出發前 3 個月完成提交，以留足基金會審查與處理時間。

A) Scholarship 申請/送件時程（倒推 T = 預計出發日）
時點	目標	你要做的事（重點）	產出物（給 sponsor 用）
T-9～T-6 個月	確立 sponsor 管線	1) 鎖定台北端 international sponsor（社/地區）與 Glasgow 端 host sponsor；2) 依 Areas of Focus 對齊提案；3) 啟動 CNA（問卷/訪談）	One-pager、CNA 表單、KPI 草案、Pilot 腳本摘要
T-6～T-4 個月	sponsor 願意「選你」並進入提案細化	1) 完成 CNA 初版洞察；2) 定義 pilot 場域與合作單位型態；3) 明確 KPI & M&E；4) 擬定 safeguarding / abuse controls	CNA 速報、KPI 表、M&E SOP、Safeguarding/Abuse policy 草案
T-4～T-3 個月	進入 Grant Center 前置	1) sponsor 整理 global grant 申請資訊（預算/合作單位/永續）；2) sponsor 在系統建立申請，並「Notify Scholarship Candidate」邀你補資料	Scholar profile（1頁）、CV、Offer/Programme 證明等（供 sponsor 上傳/引用）
T-3 個月（最晚）～T	提交與審查窗口	完成提交並配合審查、補件；同時安排出發前訓練/對外承諾（workshops/pilot 排程）	最終版專案包、里程碑、pilot 排程；（保留可稽核支出規則）

必要提醒（寫在時程表下方即可）：

個人通常不是直接向 TRF 投件；實務上由 sponsor（社/地區）建立申請並邀請你進 Grant Center 補齊 scholar profile。

Global Grant 常見門檻：總預算至少 US$30,000（由 DDF/現金/定向捐贈等組成）。

B) 12 個月實踐時程（The Hive x Anti-Fraud：產品 + 教育 + pilot）

月份以「專案啟動月 = M1」計；若你要和 scholarship 出發日對齊，通常建議 M1 ≈ T-9～T-6。

月份	Workstream	主要工作	可驗收交付物
M1–M2	CNA + 治理	台北/Glasgow 問卷（2–3週）+ 每地 6–10 場訪談；定義 taxonomy、severity、evidence gate；完成 data minimization、abuse prevention、對外溝通準則	CNA synthesis（Top risks/Top friction/Adoption constraints）、Safeguarding v1、KPI baseline 設定
M3	MVP（pilot-ready）	Report intake、Triage queue、Decision log、Audit trail；CAP message composer（模板→XML）；Education hub v1（反詐/危機溝通）	MVP demo、OpenAPI v0.1、KPI 匯出 CSV、操作手冊 v0.1
M4–M5	教材化（Supporting education）	3 條族群教案（學生/通勤/高齡或低數位素養）+ train-the-trainer；前後測題庫與滿意度量測；可及性基線（字級/對比/英文版）	Training Kit v1、Pre/Post test v1、Workshop deck v1、M&E SOP v1
M6	Pilot 準備	確認兩地 pilot 場域與合作單位；跑一次內部 tabletop；完成演練風險控管（封閉通道、不引發恐慌）	Pilot plan（場域/角色/通報鏈）、Tabletop v1（可直接跑）
M7–M8	Pilot#1（台北）	1 次 tabletop + 1 次 drill；每次都做 KPI 蒐集與 AAR；依結果修正流程與教材	台北 Pilot Report（含 KPI + AAR）、改善清單（backlog）
M9–M10	Pilot#2（Glasgow/Scotland）	同上：tabletop + drill；強化多語/在地化通報習慣；建立 host sponsor 的在地協作節點	Glasgow Pilot Report（含 KPI + AAR）、在地化差異清單
M11	擴散與永續	追加 2 個合作夥伴（目標總 ≥4）；把流程與教材做成可複製包；建立維運責任與權限分工	Partner kit（MOU 範本/Runbook/教材包）、維運手冊 v1
M12	對外交付	年度 Impact Report（KPI、案例、教材、治理）；Roadmap v2（下一年擴點與成本模型）	Year-1 Impact Report、Roadmap v2、開源發佈與版本標記
