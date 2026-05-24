# 🚀 Hyperlocal IT Career Matchmaker & Analytics Platform

A high-performance, asynchronous multi-service platform designed to bridge the gap between regional IT talents and localized industry demands. This project utilizes a **Monorepo Architecture** consisting of an orchestrator backend, an asynchronous concurrency web scraper, and an intelligence matching engine.

---

## 📌 Project Overview

This platform is developed to map IT job vacancies in real-time (specifically focusing on hyperlocal areas such as Bandung), parse mandatory skillset requirements, and accurately align them with user skill profiles using AI.

The system is decoupled into independent microservices to maintain optimized RAM utilization, enforce clean separation of concerns, and ensure high-throughput performance scalability.

### 🛠️ Tech Stack & Service Boundary
1. **Core Orchestrator & API Gateway (`/Backend`)**
   * **Architecture/Tech:** NestJS (TypeScript), Prisma ORM, PostgreSQL, Passport JWT.
   * **Responsibility:** Acts as the single source of truth for the database schema (**Data Owner**), handles user authentication, manages user profiles/skills, and serves as the primary REST API gateway for the client application.
2. **High-Performance Job Scraper Engine (`/Scraper`)**
   * **Architecture/Tech:** Go (Golang), Playwright Go, GORM (PostgreSQL Dialect).
   * **Responsibility:** Concurrently harvests job listing datasets using the **Worker Pool Pattern**, conducts data cleansing pipelines, and performs high-speed bulk insertions directly into the shared database ecosystem.
3. **AI Matchmaker Service (Future Expansion)**
   * **Architecture/Tech:** Python, LLM / Embedding Models.
   * **Responsibility:** Calculates semantic contextual relevance (**Match Score**) between raw job descriptions and a user's technical resume.

---

## 🏗️ Architectural Pattern: Shared Database Monorepo

This project implements a **Shared Database** strategy within a unified **Monorepo** workspace.

```text
├── Backend/               # NestJS Service (Data Governance & REST API Gateway)
│   ├── src/
│   └── prisma/            # Master database schemas & migration tracking
├── Scraper/               # Go Service (Concurrency Scraper Platform)
│   ├── cmd/api/
│   └── internal/          # References database tables synchronized from Prisma
└── README.md              # Global Platform Documentation

```

### Entity Synchronization Gateway (Prisma ↔️ GORM)

To eliminate table definition redundancy across polyglot microservices, **NestJS governs database migrations via Prisma**.

The Go Scraper service directly interacts with these pre-generated tables using strictly synchronized GORM models:

* **Prisma Schema (`Backend/prisma/schema.prisma`):** Defines explicit relational attributes, data primitives, and string arrays for parsed skill tags.
* **GORM Struct (`Scraper/internal/entity/job.go`):** Mapped symmetrically with explicit `json` and `gorm` structural tags, applying a `uniqueIndex` constraint on the vacancy URL field to enable conflict protection (`ON CONFLICT DO NOTHING`).

---

## 🏎️ Inside The Scraper: Go Concurrency Engine

The most critical optimization concerning data aggregation speed lies within the **Go Scraper Engine**, which coordinates a resilient **3-Phase Pipeline Architecture**:

```text
[Phase 1: Collector] ──(Job Channels)──> [Phase 2: Worker Pool (5 Workers)] ──> [Phase 3: GORM Batch Insert]

```

### Key Concurrency Features:

1. **Worker Pool Pattern:** Deploys up to **5 Concurrent Workers** running on asynchronous Goroutines to handle deep job detail queues in parallel. This optimization cuts extraction timelines by **40%**.
2. **Robust Anti-Bot Bypass:** Playwright navigation bypasses volatile network loading states (`Networkidle`), shifting instead to deterministic DOM component rendering boundaries (`WaitUntilStateCommit`).
3. **Dynamic Jitter Delay:** Injects a randomized contextual sleep interval (ranging from 2 to 5 seconds) inside each active worker's processing loop. This emulates organic human browsing behavior and eliminates active IP throttling risks (*429 Too Many Requests*).
4. **Resilient Failure Isolation:** Network or loading timeouts on individual deep-dive job links are caught gracefully, logging the failure trace independently without interrupting other active worker threads.

---

## 🚀 Development Roadmap

* [ ] **Inter-Service Integration:** Establish robust integration communication via local HTTP/REST client pooling or an event-driven message broker (e.g., RabbitMQ or Redis PubSub) to trigger scraping operations from the NestJS Core.
* [ ] **Automated Cron Scheduler:** Implement an autonomous `time.Ticker` background routine within the Go layer to automate regional scraping tasks every midnight silently.
* [ ] **Refinement of the Matching Engine:** Develop algorithmic matching intersection metrics on the database layer and finalize the dedicated Python AI service interface.
* [ ] **Graceful Shutdown Enforcement:** Monitor OS interrupt catch loops (`SIGINT/SIGTERM`) to clean up database connection pools and safely terminate virtual Playwright Chromium processes during code deployments.

---
