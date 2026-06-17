---
name: generate-qwen-md
description: Systematically explore a repository and generate a comprehensive QWEN.md context document for future interactions.
source: auto-skill
extracted_at: '2026-06-17T09:38:08.795Z'
---

# Generate a comprehensive QWEN.md

Use this skill when asked to analyze a repository and produce a `QWEN.md` (or similar project context document) that future agents can rely on.

## Procedure

1. **Initial broad exploration**
   - List the root directory contents.
   - Read the root `README.md` (or `README.txt`) if it exists. This is the fastest way to understand purpose and high-level structure.
   - Identify whether the directory is a **code project** (look for `package.json`, `go.mod`, `Cargo.toml`, `pom.xml`, `build.gradle`, `src/`, etc.) or a **non-code project** (documentation, notes, research, etc.).

2. **Iterative deep dive (up to ~10 focused reads)**
   - Based on findings, choose the next most informative files. Do not pre-plan all 10; let each discovery guide the next.
   - For code projects, prioritize in this order:
     1. Package/dependency manifest (`package.json`, `go.mod`, `requirements.txt`, etc.).
     2. Build/run configuration (`docker-compose.yml`, `Dockerfile`, `Makefile`, `angular.json`, etc.).
     3. Entrypoints (`cmd/server/main.go`, `src/main.ts`, etc.).
     4. Routing/configuration wiring (`routes.go`, `app.routes.ts`, `app.config.ts`, etc.).
     5. Environment/example files (`.env.example`).
     6. Testing/linting setup (`package.json` scripts, `eslint.config.js`, `.lintstagedrc`, test directories).
     7. Any existing project-level docs (`AGENTS.md`, older `README.md` in subpackages).
   - For non-code projects, prioritize:
     1. Index/overview files.
     2. Key documents that explain the directory's purpose.
     3. Structure/usage notes.

3. **Synthesize content by project type**

   **For a code project, produce sections for:**
   - **Project overview** — purpose, main technologies, architecture.
   - **Repository layout** — a concise tree or table of the most important directories/files and what they contain.
   - **Technologies and dependencies** — backend/frontend/dev-tool lists with brief role descriptions.
   - **Environment configuration** — how to create `.env`, required vs optional variables, redirect/cors matching caveats, generated files (e.g. `environment.auto.ts`).
   - **Building and running** — full-stack Docker command, backend-only command, frontend-only command, ports.
   - **Testing and verification** — test, lint, typecheck commands, any intentionally excluded suites (e.g. e2e).
   - **Development conventions** — architecture style (clean/hexagonal), entrypoints, route wiring, DI, required services, naming quirks, pre-commit hooks.
   - **Common gotchas** — things that trip up new contributors (e.g. Redis required at startup, legacy docs with stale paths, exact redirect URI matching, generated files not to hand-edit).
   - **Useful files to know** — a quick-reference table.

   **For a non-code project, produce sections for:**
   - **Directory overview** — what the directory is for.
   - **Key files** — the most important files and their contents.
   - **Usage** — how the contents are intended to be used.

4. **Cross-check against existing project docs**
   - If `AGENTS.md` or similar project-level instructions exist, align the generated `QWEN.md` with them.
   - Note any contradictions (e.g. older `README.md` references obsolete paths like `src/main.go` while the real entrypoint is `cmd/server/main.go`).

5. **Formatting**
   - Use well-structured Markdown: headings, tables, fenced code blocks for commands, concise bullets.
   - Keep the document scannable; future agents will read it under time pressure.

## Why this works

- Starting with the README and root listing avoids over-indexing on a single file.
- Reading iteratively lets surprising discoveries (legacy paths, generated files, required services) shape the final document.
- Separating "overview", "commands", "conventions", and "gotchas" gives future agents multiple entry points depending on what they need.

## How to apply

- When a user says "generate a QWEN.md" or "analyze this repo and create project docs", run this procedure.
- After writing the file, summarize what you produced and mention any caveats that are especially likely to trip someone up.
- If the project already has a `QWEN.md`, treat this as an update: read the existing file, then refresh stale sections while preserving its structure.
