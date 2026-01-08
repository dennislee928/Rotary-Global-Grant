# Data model (Draft)

## Tables
### reports
- id (uuid)
- created_at
- category
- severity_suggested
- area_hint (string, approx)
- time_window (string)
- description (text)
- evidence_refs (jsonb)
- reporter_contact_ref (nullable)

### triage_decisions
- id (uuid)
- report_id
- decided_at
- decision (accept/reject/needs_more_info/escalate)
- severity_final
- evidence_level
- rationale (text)
- decided_by (user_id)
- audit_hash (optional)

### alerts
- id (uuid)
- created_at
- status (draft/approved/published/withdrawn)
- cap_xml (text)
- public_message (text)
- channels (jsonb)
- approved_by (user_id)
- published_at

### training_events
- id
- date
- location
- audience
- attendance_count
- pre_avg
- post_avg
- notes

### audit_log
- id
- ts
- actor
- action
- object_type
- object_id
- diff (jsonb)
