# Proposta piano operativo - implementazione ADRs e baseline

- Fase 0: Inventario e baseline (1-2 settimane)
  - Mappa flussi critici, copertura test attuale, ADR iniziali
  - Deliverable: baseline document + scaffold ADR

- Fase 1: Config centralizzata e architettura modulare (2-3 settimane)
  - Implementare modulo config, contracts/interfaces, layout mono/multi-repo
  - Deliverables: ADR-001, ADR-002, ADR-003; config module scaffold

- Fase 2: Build/CI/CD/DX (2-3 settimane)
  - Unificare scripts, TS strict, alias, logging, security policy
  - Deliverables: package.json scripts, tsconfig paths, lint/format, logging lib

- Fase 3: Qualità e osservabilità (3-5 settimane)
  - Copertura test target, test di integrazione, CI/CD checks
  - Deliverables: tests, coverage reports, monitoring hooks

- Fase 4: Documentazione e governance (ongoing)
  - ADRs, README, CONTRIBUTING, update policy

- KPI
  - Copertura test, CI stability, baseline/config validation, governance docs.
