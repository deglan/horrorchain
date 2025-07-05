# 🧠 Blockchain Horror Game (Go + Ebiten)

An experimental 2D top-down psychological horror game that teaches the fundamentals of blockchain through immersive gameplay and atmosphere.

> Engine: [Ebiten (Go)](https://ebiten.org)  
> Theme: Psychological horror + technical learning  
> Format: Open-source educational indie project

---

## 📁 Project Structure

├── assets/ → Game assets (graphics, audio, fonts)
│ ├── player/ → Player sprites
│ ├── tileset/ → Environment tiles (map blocks)
│ ├── audio/ → Ambient and sound effects
│ └── fonts/ → Pixel or glitch fonts
│
├── constants/ → Constants
│
├── blocks/ → Blockchain logic and data structures
│ ├── block.go → Block struct (index, hash, data…)
│ └── chain.go → Blockchain validation, linking
│
├── scenes/ → Game scenes (start screen, gameplay)
│ ├── start.go → Start / intro screen
│ └── game.go → Main game logic
│
├── engine/ → Core engine: input, movement, rendering
│ ├── player.go → Player controller
│ └── map.go → Map renderer and loader
│
├── main.go → Main game loop and entry point
└── Makefile → Common tasks: run, build, lint, clean


---

## 🚀 Running the Game

Using Make:
```bash
make run
Or directly with Go: go run main.go

