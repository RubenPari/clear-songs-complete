# Operational Plan - ADRs & Baseline (READ/BUILD)

Obiettivo: definire un piano operativo task-level per centralizzare configurazione, modularizzare l'architettura, e migliorare CI/CD/DX, QA e observability.

Scope
- Config centralizzata, contracts/interfaces, ADRs iniziali
- Build/CI/CD/DX, linting, TS strict, alias, logging
- QA, test coverage, security, observability
- Documentazione e governance

Phases, Deliverables e KPI
- Fase 0 — Inventario e baseline (0-2w)
  - Tasks: mappa flussi critici; baseline test coverage; scaffolding ADR
  - Deliverable: baseline.md; ADR scaffolds
  - KPI: Flussi chiave identificati; baseline test coverage definita
- Fase 1-2 — Config centralized + Architettura modulare (2-4w)
  - Tasks: config module (env, validation, defaults); contracts/interfaces; ADR updates
  - Deliverables: config module scaffold; ADR-001/ADR-002 aggiornate
  - KPI: validazione config in startup; contracts in place
- Fase 3-5 — Build/CI/CD/DX (4-6w)
  - Tasks: script unificati (lint/test/build/start) in package.json; TS strict; path alias; logging
  - Deliverables: package.json scripts; tsconfig paths; lint config; logging lib
  - KPI: CI stabile; lint/strict attivi; copertura in crescita
- Fase 6-8 — QA, osservabilità e sicurezza (6-10w)
  - Tasks: aumentare test; integrazione test di integrazione; monitoring hooks
  - Deliverables: suite di test; reports di coverage; monitoring hooks
  - KPI: copertura target; riduzione fail critici
- Fase 9 — Governance e docs (ongoing)
  - Tasks: ADRs aggiornate; CONTRIBUTING; policy dipendenze
  - Deliverables: README/CONTRIBUTING; PRIORITIZATION aggiornata
  - KPI: ADRs allineate; onboarding facilitato

Deliverables generali
- ADRs (001-004) aggiornate o create
- Config module scaffold
- Script di build/lint/test/documentazione aggiornata
- Piani di test, logging e observability implementati

Rischi principali e mitigazioni
- Resistenza al cambiamento: formazione e walkthrough
- Integrazione graduale: pipeline CI/CD con feature flag di migrazione
- Dipendenze esterne: auditing periodico delle dipendenze

Next steps
- Confermare stack e formato ADR (Markdown/ADR-first)
- Definire ownership per Fase e owner per task
- Preparare un primo runbook di setup dev
