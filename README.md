# 🎮 F2P Game Explorer

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com)
[![SSE](https://img.shields.io/badge/SSE-Streaming-FF6B6B?style=for-the-badge&logo=lightning&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
[![Tailwind](https://img.shields.io/badge/Tailwind-CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)](https://tailwindcss.com)
[![Render](https://img.shields.io/badge/Render-Deployed-46E3B7?style=for-the-badge&logo=render&logoColor=white)](https://render.com)

Real-time free-to-play game discovery stream using Go, Server-Sent Events (SSE), and FreeToGame API. Explore 500+ F2P games across PC and browser platforms! 🕹️

## ✨ Features

- 🎮 **500+ Games** — Discover free-to-play games from FreeToGame database
- 🔴 **Real-time SSE** — Server-Sent Events for live streaming
- 🎨 **Gaming UI** — Cyberpunk-themed interface with neon effects
- 🏷️ **Rich Metadata** — Genre, platform, publisher, release date
- 🔗 **Direct Links** — Play games directly from the explorer
- 📱 **Responsive** — Works on all screen sizes
- ⚡ **Fast** — Built with Gin high-performance framework

## 🚀 Quick Start

Clone the repository:

```bash
git clone https://github.com/smart-developer1791/go-game-explorer
cd go-game-explorer
```

Initialize dependencies and run:

```bash
go mod tidy
go run .
```

Open browser at **http://localhost:8080** 🎮

## 🎯 API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Main UI interface |
| `GET /stream` | SSE game stream |
| `GET /stats` | API statistics |

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.21+ |
| Framework | Gin |
| Streaming | Server-Sent Events |
| API | FreeToGame API |
| Styling | Tailwind CSS |
| Fonts | Orbitron, Rajdhani |

## 📊 Game Categories

| Genre | Examples |
|-------|----------|
| MMORPG | World of Warcraft, Lost Ark |
| Shooter | Valorant, Apex Legends |
| MOBA | League of Legends, Dota 2 |
| Battle Royale | Fortnite, PUBG |
| Card Game | Hearthstone, Legends of Runeterra |
| Strategy | StarCraft II, Age of Empires Online |

## 🎨 UI Features

- 💜 **Neon Cyberpunk Theme** — Purple and cyan gradient design
- ✨ **Animated Cards** — Smooth slide-in animations
- 🌟 **Hover Effects** — Interactive game cards
- 🏷️ **Smart Badges** — Platform and genre indicators
- 🎮 **Gaming Fonts** — Orbitron for headings

## 📁 Project Structure

```text
go-game-explorer/
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── render.yaml      # Render deployment config
├── .gitignore       # Git ignore rules
└── README.md        # Documentation
```

## 🌐 Data Source

This project uses the [FreeToGame API](https://www.freetogame.com/api-doc):

- 🆓 **Free to use** — No API key required
- 📊 **500+ games** — Comprehensive F2P database
- 🔄 **Regular updates** — New games added frequently
- 📋 **Rich data** — Thumbnails, descriptions, metadata

## ⚙️ Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |

## 🙏 Acknowledgments

- [FreeToGame](https://www.freetogame.com/) — Game database API
- [Gin](https://gin-gonic.com/) — Web framework
- [Tailwind CSS](https://tailwindcss.com/) — Styling

---

## Deploy in 10 seconds

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)
