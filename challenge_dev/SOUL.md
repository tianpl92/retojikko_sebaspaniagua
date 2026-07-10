# Public Calls Portal — Executive and Technical Delivery Report

## 1. Project construction information

The Public Calls Portal was delivered through a structured, spec-driven workflow managed with Hermes. Development began with foundational planning and database definition, then progressed through backend skeleton creation, REST endpoint implementation, frontend shell development, backend/frontend integration, iterative bug fixing, and final documentation hardening.

The project was intentionally organized as a sequence of numbered specifications in `specs_project/`, each one describing a bounded objective, the target files, expected results, and verification requirements. That approach made the delivery process auditable and reduced drift between the intended product behavior and the implemented code.

At a high level, the work evolved in four major stages:

1. Data model and persistence design in PostgreSQL.
2. Go backend implementation with authentication, proposal browsing, persistence, and session handling.
3. Static frontend implementation in HTML/CSS/JavaScript.
4. Integration, stabilization, and documentation/reporting.

The result is a complete web application for browsing public calls from datos.gov.co, saving proposals, authenticating users via JWT, and maintaining a local persistent record of user activity.

## 2. Architecture

The final solution follows a layered, full-stack architecture with clear separation of concerns.

### 2.1 Backend architecture

The backend is implemented in Go and organized into domain, repository, service, handler, router, and database layers:

- `backend/cmd/api/` contains the application entrypoint.
- `backend/internal/config/` centralizes runtime configuration.
- `backend/internal/database/` manages PostgreSQL connection and schema validation/bootstrap.
- `backend/internal/domain/` defines core business entities such as users, public proposals, saved proposals, and user sessions.
- `backend/internal/repository/` contains repository interfaces, mock implementations, and PostgreSQL implementations.
- `backend/internal/service/` contains the business logic for authentication, proposal retrieval, saving, and external integration.
- `backend/internal/handler/` exposes HTTP handlers and middleware for the REST API.
- `backend/internal/router/` wires routes, middleware, and CORS behavior.

This architecture allows the application to run in two modes:

- production-style mode using PostgreSQL-backed repositories;
- development/testing mode using in-memory mock repositories when the database is unavailable.

### 2.2 Persistence architecture

The data layer uses PostgreSQL as the system of record. The repository layer abstracts persistence so that the business logic is not coupled directly to SQL or storage details. The schema includes the principal business objects and session tracking, allowing the backend to authenticate users, persist proposal data, and invalidate sessions during logout.

A schema validation/bootstrap process runs at backend startup to ensure the database is ready. If the schema is already present, the system skips re-creation. If it is missing, the project applies the SQL schema from `database/public_calls_database.sql`.

### 2.3 Frontend architecture

The frontend is a static client implemented with HTML5, CSS, and vanilla JavaScript. It is intentionally lightweight and served separately with `python3 -m http.server`.

Frontend pages include:

- login page
- dashboard for browsing public calls
- saved proposals page
- user creation page

The frontend communicates with the backend via REST/JSON using `fetch`, and it handles local session state using browser storage. The UI is designed in Spanish and uses the project color palette defined in the repository guidance.

### 2.4 Integration architecture

The application integrates an external open-data source from datos.gov.co, which feeds the public proposals browsing workflow. The backend acts as the orchestration and normalization layer between the external data source, the local persistence model, and the user-facing frontend.

This design keeps the frontend simple and keeps external data normalization centralized in the backend services.

## 3. Hermes usage

Hermes was used as the workflow and coordination layer for the entire project lifecycle.

### 3.1 Planning and execution discipline

The implementation was organized around a stored Hermes plan in `.hermes/plans/2026-07-02_161801-public-calls-portal.md`. That plan captured the project goal, the architecture direction, the milestone sequence, and the verification checklist. It served as the main control document for the build process.

### 3.2 Spec-driven delivery

The project was advanced through numbered spec documents in `specs_project/`. Hermes helped keep each stage bounded and traceable, preventing uncontrolled implementation jumps. This was especially valuable when the work shifted from core MVP delivery into refinement and stabilization tasks.

### 3.3 Agent and subagent workflow

The project used a multi-agent approach for reasoning, implementation, and review.

- The principal Hermes agent coordinated planning, document review, and progress control.
- Planning-oriented tasks were handled with a lighter orchestration model suited for structuring work and producing specifications.
- Code-generation tasks, especially Go backend implementation and frontend UI work, were handled by a stronger implementation model through subagents.

This division of labor improved speed and quality:

- the main agent stayed focused on orchestration and correctness;
- subagents produced focused outputs for backend and frontend concerns;
- the workflow reduced context overload and kept each task narrowly scoped.

### 3.4 Practical value

Hermes provided:

- controlled decomposition of work;
- traceability across specs, commits, and outcomes;
- a repeatable plan → execute → verify pattern;
- documentation continuity from the first architecture decisions to the final report.

## 4. Why I use ?

The project used multiple models according to the nature of the task and the expected complexity.

### 4.1 DeepSeek v4 Flash

DeepSeek v4 Flash was used for planning, document generation, instruction writing, and orchestration tasks where speed and structured reasoning were more important than heavy code synthesis.

It was especially suitable for:

- preparing specs;
- restructuring documentation;
- generating implementation plans;
- coordinating multi-step work.

### 4.2 DeepSeek v4 Pro

DeepSeek v4 Pro was used for implementation-heavy tasks, especially Go backend code and frontend JavaScript/HTML work. It was the preferred model for generating concrete code and handling more detailed engineering output.

It was especially useful for:

- Go handlers, services, repositories, and tests;
- frontend page logic and interaction code;
- code review and integration refinement.

### 4.3 GPT-5.5 and GPT-5.4

GPT-5.5 and GPT-5.4 were used as additional support models during the project lifecycle, particularly when coverage, experimentation, or alternate reasoning paths were useful.

The combination of models allowed the workflow to balance:

- speed for planning,
- quality for implementation,
- flexibility for iterative refinement.

## 5. Blocks

The main technical blockers were integration mismatches, schema/bootstrap concerns, and payload-shape differences between systems.

### 5.1 Main blocker categories

1. Data-shape mismatches between frontend requests and backend expectations.
2. Response-shape inconsistencies when integrating saved proposals and external data.
3. Session and logout handling, including token invalidation and cleanup.
4. Database bootstrap and schema validation during backend startup.
5. External API normalization for datos.gov.co payloads.

### 5.2 How the blockers were resolved

These issues were resolved through the spec-driven workflow rather than ad hoc patching. Each spec targeted one bounded problem and produced a verifiable change set.

Key technical fixes included:

- mapping `id_del_proceso` to the backend `id` field before save requests;
- normalizing `urlproceso` because datos.gov.co returns it as an object rather than a plain string;
- aligning saved proposals response structures with the frontend modal and listing requirements;
- implementing logout so JWT sessions can be invalidated server-side and cleared client-side;
- adding schema validation/bootstrap logic so the backend can initialize safely.

### 5.3 Most relevant commits

The following commits best represent the project’s evolution and the major milestones of the delivered architecture:

- `a6035d9` — feat: challenge dev database schema creation
- `18d3ba1` — feat: database schema v2 with SECOP II column comparison
- `95256fb` — refactor SQL schema
- `d68b6b1` — feat: add backend skeleton with layered architecture
- `5aa5167` — feat: add JWT auth, user endpoints, and real PostgreSQL repos
- `4e689f7` — feat: add specs 1-9, integration docs, JWT backend
- `8e5468b` — build integration with external API and initial API contract
- `215fb84` — build logic for saved proposals, user info, user modify, and saved proposals retrieval
- `23a5bf9` — dashboard skeleton creation and login integration
- `b032147` — integrate dashboard and saved proposals pages
- `8de49da` — bug fixes about saved proposals
- `a57b094` — construction challenge finished, adding README.md and SOUL.md
- `8e185a7` — update repository URL in SOUL.md

These commits tell the story of the project from schema foundation to application stabilization and final documentation.

## 6. Suggestions

The project is complete and functionally coherent, but there are several recommendations for future maintenance and expansion.

### 6.1 Technical recommendations

- Keep the spec-driven workflow for future enhancements.
- Maintain the repository/service/handler separation to preserve testability.
- Continue validating payload shapes early when consuming external APIs.
- Prefer explicit schema bootstrap and migration control for database-dependent work.
- Expand automated tests whenever new endpoints or data shapes are introduced.

### 6.2 Process recommendations

- Preserve the plan → execute → verify discipline for future releases.
- Continue using subagents for isolated backend/frontend implementation tasks.
- Capture architecture decisions in a single living report to keep the project explainable over time.

### 6.3 Operational recommendations

- Keep the backend startup validation behavior documented.
- Maintain the frontend/server separation for easier local development.
- Review commit history alongside specs when onboarding new contributors.

## 7. Repository github

Repository URL:

https://github.com/tianpl92/retojikko_sebaspaniagua/tree/main/challenge_dev

### 7.1 Documentation and reference paths

The following references capture the main project documentation and the specification trail that guided the delivery process:

- Main project README:
  https://github.com/tianpl92/retojikko_sebaspaniagua/blob/main/challenge_dev/README.md
- Backend README:
  https://github.com/tianpl92/retojikko_sebaspaniagua/blob/main/challenge_dev/backend/README.md
- Spec documents directory:
  https://github.com/tianpl92/retojikko_sebaspaniagua/tree/main/challenge_dev/specs_project/

These paths are useful as the canonical entry points for understanding the repository structure, local execution instructions, and the specification sequence used during implementation.

---

## Appendix: executive summary

The Public Calls Portal is a complete, spec-driven web application built with a Go backend, PostgreSQL persistence, and a static HTML/CSS/JavaScript frontend. The architecture is layered, testable, and intentionally simple to operate. Hermes was used to plan, coordinate, and document the work, while subagents supported focused generation and refinement tasks. The result is a mature delivery with clear technical boundaries, a stable integration story, and a commit history that reflects each major implementation phase.