# Renbun (錬文 - Sentence Alchemy)
A high-velocity, distributed Japanese literature miner and Discord bot built in Go.

🎯 **Project Overview & Motivation**
Renbun was engineered to automate the extraction of native Japanese sentences from the live web for immersion and spaced repetition (SRS) pipelines.
Moving beyond sequential Python scraping, this engine demonstrates production-grade backend architecture by utilizing Go's M:N scheduler, concurrent worker pools, and persistent WebSocket connections to navigate and refine raw HTML into structured linguistic data.

⚡ **System Architecture**
`[Discord UI]` <-> `[WebSocket Gateway]` <-> `[Goroutine Worker Pool]` <-> `[Target DOMs]`

* **The Frontend:** Headless command execution via the Discord API (WebSocket).
* **The Orchestrator:** A Go-based worker pool that prevents resource exhaustion while concurrently processing job queues.
* **The Extraction:** DOM parsing utilizing `goquery` with automatic Shift_JIS/EUC-JP character encoding decoding and Furigana stripping.
* **The Filtering Pipeline:** Length boundary checks (`15 <= len <= 60 runes`), keyword matching, and random array shuffling.

🚀 **Commands**
* `!ping` — Verifies engine status and latency over the WebSocket gateway.
* `!mine <keyword>` — Triggers the distributed worker pool across `target.txt` sources, returning 3 refined, highlighted immersion sentences in a rich embed.

🛠️ **Development Roadmap**
* [x] Phase 0: Infrastructure & Authentication provisioning.
* [x] Phase 1: WebSocket connection & Discord Event Handlers.
* [x] Phase 2: Core DOM Extraction Logic & Shift_JIS translation.
* [x] Phase 3: Distributed Worker Pool Implementation.
* [x] Phase 4: Integration, Mathematical Filtering, and Rich UI.
