# Renbun (錬文 - Sentence Alchemy)
A high-velocity, distributed Japanese literature miner and Discord bot built in Go.

🎯 **Project Overview & Motivation**
Renbun was engineered to automate the extraction of native Japanese sentences from the live web for immersion and spaced repetition (SRS) pipelines.
Moving beyond sequential Python scraping, this engine demonstrates production-grade backend architecture by utilizing Go's M:N scheduler, concurrent worker pools, and persistent WebSocket connections to navigate and refine raw HTML into structured linguistic data.

⚡ **System Architecture**
`[Discord UI]` <-> `[WebSocket Gateway]` <-> `[Goroutine Worker Pool]` <-> `[Target DOMs]`

* **The Frontend:** Headless command execution via the Discord API (WebSocket).
* **The Orchestrator:** A Go-based worker pool that prevents resource exhaustion while concurrently processing job queues.
* **The Extraction:** DOM parsing utilizing `goquery` to isolate native Japanese structures from HTML noise.

🛠️ **Development Roadmap**
* [x] Phase 0: Infrastructure & Authentication provisioning.
* [ ] Phase 1: WebSocket connection & Discord Event Handlers.
* [ ] Phase 2: Core DOM Extraction Logic.
* [ ] Phase 3: Distributed Worker Pool Implementation.
