# Public Calls Portal — Technical Status Report

## 1. Project construction information

The project was built through a spec-driven workflow guided by the Hermes plan in `.hermes/plans/2026-07-02_161801-public-calls-portal.md`. The work started with database design, continued with backend skeleton and API implementation, then frontend integration, and finally refinement, bug fixing, and documentation updates. Each numbered spec in `specs_project/` defined the context, target files, and expected result for a controlled delivery sequence.

## 2. Architecture

The solution is split into three main layers:

- `backend/`: Go HTTP API with handlers, services, repositories, JWT middleware, and startup schema validation.
- `frontend/`: static HTML/CSS/JavaScript client served with `python3 -m http.server`.
- `database/`: PostgreSQL schema script used for local provisioning.

Core technologies:

- Go + pgx for backend and database access
- PostgreSQL for persistence
- HTML5/CSS/JavaScript for the UI
- JWT for authentication
- REST/JSON for API communication

## 3. Hermes usage

Hermes was used as the project orchestration layer to read the plan, follow the numbered specs, and keep implementation aligned with the expected flow. The main references were the database specs and the skeleton backend service specs, which guided the schema, repository structure, and startup behavior. Hermes also helped track the iterative fixes across later specs and maintain a short technical report of the delivered work.

## 4. Why I use ?

The project used multiple models depending on the task type:

### 4.1 DeepSeek v4 Flash
I select deepseek v4 flash model for running plans, make documention nad write instructions, and by the cost.

### 4.2 DeepSeek v4 Pro
I select deepseek pro model for coding and execute complex taskes and by the cost.

### 4.3 GPT-5.5 and GPT-5.4
Additional I use gpt5.5 and gpt5.4, starting project by free license until limit ussage.

## 5. Blocks

The main blockers were integration mismatches and data-shape issues during development. These were resolved through the executed specs, including the database specification set, the backend skeleton service spec, the saved proposals fixes, the logout feature, and the startup schema validation work. The main technical issues fixed were response-shape alignment, `urlproceso` input normalization, session invalidation, and schema bootstrap at backend startup.

## 6. Suggestions 

To understand and make time to keep learning about the skills Hermes handles and how to integrate it with other platforms.

## 7. Repository github 

https://github.com/tianpl92/retojikko_sebaspaniagua/tree/main
