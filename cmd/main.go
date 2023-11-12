package main

import (
	"embed"
	"fmt"
	"gsnake/pkg/env"
	"gsnake/pkg/game"
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

type config struct {
	host string
	port uint16
}

//go:embed templates
var templateFiles embed.FS

//go:embed static
var staticFiles embed.FS

func main() {
	var cfg config
	cfg.host = env.GetEnv("HOST", "localhost")
	cfg.port = uint16(env.GetEnvInt("PORT", 8000))

	// create server
	s := server{
		nil,
		slog.Default(),
		make(map[string]*game.Game),
	}
	s.t = template.Must(template.New("snake").Funcs(template.FuncMap{"IsGameOver": s.isGameOver}).ParseFS(templateFiles, "templates/*"))

	// set routes
	http.HandleFunc("/", s.index)
	http.HandleFunc("/game", s.render)
	http.HandleFunc("/direction", s.direction)
	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	// start
	s.log.Info(fmt.Sprintf("listening on %s:%d", cfg.host, cfg.port))
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.host, cfg.port), nil))
}
