package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Game struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Thumbnail        string `json:"thumbnail"`
	ShortDescription string `json:"short_description"`
	GameURL          string `json:"game_url"`
	Genre            string `json:"genre"`
	Platform         string `json:"platform"`
	Publisher        string `json:"publisher"`
	Developer        string `json:"developer"`
	ReleaseDate      string `json:"release_date"`
}

type GameStore struct {
	games []Game
	mu    sync.RWMutex
}

var store = &GameStore{}

func (s *GameStore) SetGames(games []Game) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games = games
}

func (s *GameStore) GetRandom() (Game, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.games) == 0 {
		return Game{}, false
	}
	return s.games[rand.Intn(len(s.games))], true
}

func (s *GameStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.games)
}

func fetchGames() error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://www.freetogame.com/api/games")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var games []Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return err
	}

	store.SetGames(games)
	log.Printf("✅ Loaded %d games from FreeToGame API", len(games))
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := fetchGames(); err != nil {
		log.Printf("⚠️ Warning: Could not fetch games: %v", err)
	}

	// Refresh games periodically
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			if err := fetchGames(); err != nil {
				log.Printf("⚠️ Refresh failed: %v", err)
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/", indexHandler)
	r.GET("/stream", streamHandler)
	r.GET("/stats", statsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🎮 Game Explorer running on http://localhost:%s", port)
	r.Run(":" + port)
}

func indexHandler(c *gin.Context) {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	c.Header("Content-Type", "text/html")
	tmpl.Execute(c.Writer, nil)
}

func statsHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_games": store.Count(),
		"status":      "online",
	})
}

func streamHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(500, gin.H{"error": "SSE not supported"})
		return
	}

	sendGame := func() {
		game, ok := store.GetRandom()
		if !ok {
			fmt.Fprintf(c.Writer, "event: error\ndata: {\"message\":\"No games available\"}\n\n")
			flusher.Flush()
			return
		}
		data, _ := json.Marshal(game)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	}

	// Send first game immediately
	sendGame()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			sendGame()
		}
	}
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🎮 F2P Game Explorer</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;700;900&family=Rajdhani:wght@400;500;700&display=swap');
        
        body {
            font-family: 'Rajdhani', sans-serif;
            background: linear-gradient(135deg, #0f0c29 0%, #302b63 50%, #24243e 100%);
        }
        
        .font-gaming { font-family: 'Orbitron', monospace; }
        
        .neon-border {
            box-shadow: 0 0 5px #00ffff, 0 0 10px #00ffff, inset 0 0 5px rgba(0,255,255,0.1);
            border: 1px solid #00ffff;
        }
        
        .neon-text {
            text-shadow: 0 0 10px #00ffff, 0 0 20px #00ffff, 0 0 30px #00ffff;
        }
        
        .neon-purple {
            box-shadow: 0 0 5px #ff00ff, 0 0 10px #ff00ff;
            border: 1px solid #ff00ff;
        }
        
        .game-card {
            background: rgba(15, 12, 41, 0.9);
            backdrop-filter: blur(10px);
            transition: all 0.3s ease;
            animation: slideIn 0.5s ease-out;
        }
        
        .game-card:hover {
            transform: translateY(-5px) scale(1.02);
            box-shadow: 0 10px 30px rgba(0, 255, 255, 0.3);
        }
        
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(30px) scale(0.9);
            }
            to {
                opacity: 1;
                transform: translateY(0) scale(1);
            }
        }
        
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        
        .pulse { animation: pulse 2s infinite; }
        
        @keyframes float {
            0%, 100% { transform: translateY(0); }
            50% { transform: translateY(-10px); }
        }
        
        .float { animation: float 3s ease-in-out infinite; }
        
        .gradient-text {
            background: linear-gradient(90deg, #00ffff, #ff00ff, #00ffff);
            background-size: 200% auto;
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            animation: shine 3s linear infinite;
        }
        
        @keyframes shine {
            to { background-position: 200% center; }
        }
        
        .genre-badge {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        
        .platform-pc {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }
        
        .platform-browser {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }
        
        .btn-neon {
            position: relative;
            overflow: hidden;
            transition: all 0.3s;
        }
        
        .btn-neon::before {
            content: '';
            position: absolute;
            top: 0;
            left: -100%;
            width: 100%;
            height: 100%;
            background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
            transition: left 0.5s;
        }
        
        .btn-neon:hover::before {
            left: 100%;
        }
        
        .scrollbar-gaming::-webkit-scrollbar {
            width: 8px;
        }
        
        .scrollbar-gaming::-webkit-scrollbar-track {
            background: rgba(0,0,0,0.3);
            border-radius: 4px;
        }
        
        .scrollbar-gaming::-webkit-scrollbar-thumb {
            background: linear-gradient(180deg, #00ffff, #ff00ff);
            border-radius: 4px;
        }
    </style>
</head>
<body class="min-h-screen text-white scrollbar-gaming">
    <div class="container mx-auto px-4 py-8">
        <!-- Header -->
        <header class="text-center mb-10">
            <div class="float inline-block mb-4">
                <span class="text-7xl">🎮</span>
            </div>
            <h1 class="font-gaming text-4xl md:text-6xl font-black gradient-text mb-4">
                F2P GAME EXPLORER
            </h1>
            <p class="text-cyan-300 text-lg md:text-xl max-w-2xl mx-auto">
                Real-time free-to-play game discovery stream powered by 
                <span class="text-pink-400 font-bold">FreeToGame API</span>
            </p>
            <div class="mt-4 flex items-center justify-center gap-4 text-sm">
                <span class="flex items-center gap-2 bg-green-500/20 px-3 py-1 rounded-full border border-green-500/50">
                    <span class="w-2 h-2 bg-green-400 rounded-full pulse"></span>
                    <span id="status">Connecting...</span>
                </span>
                <span class="bg-purple-500/20 px-3 py-1 rounded-full border border-purple-500/50">
                    📊 <span id="discovered">0</span> discovered
                </span>
                <span class="bg-cyan-500/20 px-3 py-1 rounded-full border border-cyan-500/50">
                    🗃️ <span id="totalGames">500+</span> games
                </span>
            </div>
        </header>

        <!-- Controls -->
        <div class="flex flex-wrap justify-center gap-4 mb-10">
            <button id="startBtn" onclick="startStream()" 
                class="btn-neon neon-border bg-cyan-500/20 hover:bg-cyan-500/40 px-8 py-3 rounded-lg font-gaming font-bold tracking-wider transition-all">
                ▶️ START STREAM
            </button>
            <button id="stopBtn" onclick="stopStream()" 
                class="btn-neon neon-purple bg-pink-500/20 hover:bg-pink-500/40 px-8 py-3 rounded-lg font-gaming font-bold tracking-wider transition-all" disabled>
                ⏹️ STOP
            </button>
            <button onclick="clearGames()" 
                class="btn-neon bg-gray-500/20 hover:bg-gray-500/40 px-8 py-3 rounded-lg font-gaming font-bold tracking-wider transition-all border border-gray-500">
                🗑️ CLEAR
            </button>
        </div>

        <!-- Games Grid -->
        <div id="gamesGrid" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <!-- Games will appear here -->
        </div>

        <!-- Empty State -->
        <div id="emptyState" class="text-center py-20">
            <div class="text-8xl mb-6 opacity-50">🕹️</div>
            <h2 class="font-gaming text-2xl text-gray-400 mb-4">No Games Yet</h2>
            <p class="text-gray-500">Click START STREAM to discover free-to-play games</p>
        </div>
    </div>

    <!-- Footer -->
    <footer class="text-center py-8 text-gray-500 text-sm">
        <p>Powered by <a href="https://www.freetogame.com/api-doc" target="_blank" class="text-cyan-400 hover:text-cyan-300">FreeToGame API</a></p>
        <p class="mt-2">Built with Go + Gin + SSE + Tailwind CSS</p>
    </footer>

    <script>
        let eventSource = null;
        let discoveredCount = 0;
        let discoveredGames = new Set();

        // Fetch stats on load
        fetch('/stats')
            .then(r => r.json())
            .then(data => {
                document.getElementById('totalGames').textContent = data.total_games || '500+';
            })
            .catch(() => {});

        function startStream() {
            if (eventSource) return;

            eventSource = new EventSource('/stream');
            
            document.getElementById('status').textContent = 'Live Streaming';
            document.getElementById('startBtn').disabled = true;
            document.getElementById('stopBtn').disabled = false;
            document.getElementById('emptyState').classList.add('hidden');

            eventSource.onmessage = (event) => {
                try {
                    const game = JSON.parse(event.data);
                    
                    // Track unique games
                    if (!discoveredGames.has(game.id)) {
                        discoveredGames.add(game.id);
                        discoveredCount++;
                        document.getElementById('discovered').textContent = discoveredCount;
                    }
                    
                    addGameCard(game);
                } catch (e) {
                    console.error('Parse error:', e);
                }
            };

            eventSource.onerror = () => {
                document.getElementById('status').textContent = 'Reconnecting...';
                setTimeout(() => {
                    if (eventSource) {
                        stopStream();
                        startStream();
                    }
                }, 3000);
            };
        }

        function stopStream() {
            if (eventSource) {
                eventSource.close();
                eventSource = null;
            }
            document.getElementById('status').textContent = 'Stopped';
            document.getElementById('startBtn').disabled = false;
            document.getElementById('stopBtn').disabled = true;
        }

        function clearGames() {
            document.getElementById('gamesGrid').innerHTML = '';
            discoveredCount = 0;
            discoveredGames.clear();
            document.getElementById('discovered').textContent = '0';
            document.getElementById('emptyState').classList.remove('hidden');
        }

        function getPlatformClass(platform) {
            if (platform.toLowerCase().includes('pc')) return 'platform-pc';
            return 'platform-browser';
        }

        function getGenreEmoji(genre) {
            const emojis = {
                'mmorpg': '⚔️',
                'shooter': '🔫',
                'strategy': '🧠',
                'moba': '🏆',
                'racing': '🏎️',
                'sports': '⚽',
                'social': '👥',
                'sandbox': '🏗️',
                'survival': '🏕️',
                'card': '🃏',
                'battle-royale': '🎯',
                'fantasy': '🧙',
                'sci-fi': '🚀',
                'fighting': '🥊',
                'action': '💥',
                'adventure': '🗺️',
                'puzzle': '🧩'
            };
            const key = genre.toLowerCase().replace(' ', '-');
            return emojis[key] || '🎮';
        }

        function addGameCard(game) {
            const grid = document.getElementById('gamesGrid');
            const card = document.createElement('div');
            card.className = 'game-card rounded-xl overflow-hidden neon-border';
            
            card.innerHTML = ` + "`" + `
                <div class="relative">
                    <img src="${game.thumbnail}" alt="${game.title}" 
                         class="w-full h-48 object-cover"
                         onerror="this.src='https://via.placeholder.com/365x206/1a1a2e/00ffff?text=No+Image'">
                    <div class="absolute top-3 right-3">
                        <span class="${getPlatformClass(game.platform)} px-3 py-1 rounded-full text-xs font-bold text-white shadow-lg">
                            ${game.platform.includes('PC') ? '💻 PC' : '🌐 Browser'}
                        </span>
                    </div>
                    <div class="absolute bottom-0 left-0 right-0 h-20 bg-gradient-to-t from-[#0f0c29] to-transparent"></div>
                </div>
                <div class="p-5">
                    <div class="flex items-start justify-between mb-3">
                        <h3 class="font-gaming font-bold text-lg text-cyan-300 leading-tight flex-1 mr-2">${game.title}</h3>
                        <span class="genre-badge px-2 py-1 rounded text-xs font-bold whitespace-nowrap">
                            ${getGenreEmoji(game.genre)} ${game.genre}
                        </span>
                    </div>
                    <p class="text-gray-400 text-sm mb-4 line-clamp-2 leading-relaxed">${game.short_description}</p>
                    <div class="flex items-center justify-between text-xs text-gray-500 mb-4">
                        <span>🏢 ${game.publisher}</span>
                        <span>📅 ${game.release_date}</span>
                    </div>
                    <a href="${game.game_url}" target="_blank" 
                       class="block w-full text-center bg-gradient-to-r from-cyan-500 to-purple-500 hover:from-cyan-400 hover:to-purple-400 py-2 rounded-lg font-gaming font-bold text-sm transition-all hover:shadow-lg hover:shadow-cyan-500/25">
                        🎮 PLAY FREE
                    </a>
                </div>
            ` + "`" + `;
            
            grid.insertBefore(card, grid.firstChild);

            // Keep max 20 games
            while (grid.children.length > 20) {
                grid.removeChild(grid.lastChild);
            }
        }
    </script>
</body>
</html>`
