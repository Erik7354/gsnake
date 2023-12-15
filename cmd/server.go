package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gsnake/pkg/game"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type keyCode = string

const (
	arrowUpKey    keyCode = "ArrowUp"
	arrowDownKey  keyCode = "ArrowDown"
	arrowLeftKey  keyCode = "ArrowLeft"
	arrowRightKey keyCode = "ArrowRight"
)

type server struct {
	t     *template.Template
	log   *slog.Logger
	games map[string]*game.Game
}

func (s server) session(r *http.Request) (*game.Game, error) {
	cookie, err := r.Cookie("gSnakeId")
	if err != nil || cookie == nil {
		return nil, err
	}

	g, ok := s.games[cookie.Value]
	if !ok {
		return nil, errors.New("game not found")
	}

	return g, nil
}

func (s server) isGameOver(gameID string) bool {
	return s.games[gameID].GameOver
}

// ######################################################################## //
// ### Handler
// ######################################################################## //

func (s server) index(w http.ResponseWriter, r *http.Request) {
	// n
	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil || n < 5 {
		n = 5
	}
	// refresh rate
	rr := "1000ms"
	rrd, err := time.ParseDuration(r.URL.Query().Get("rr"))
	if err == nil {
		rr = fmt.Sprintf("%dms", rrd.Milliseconds())
	}

	// set render cookie
	var b = make([]byte, 4)
	_, err = rand.Read(b)
	if err != nil {
		http.Redirect(w, r, "/static/500-meme.avif", http.StatusFound)
		return
	}

	gid := hex.EncodeToString(b)
	s.games[gid] = game.New(n)

	http.SetCookie(w, &http.Cookie{
		Name:  "gSnakeId",
		Value: gid,
	})

	_ = s.t.ExecuteTemplate(w, "index", struct {
		N       int
		RR      string
		Session game.Game
	}{
		N:       n,
		RR:      rr,
		Session: *s.games[gid],
	})

	s.log.Info(fmt.Sprintf("[%s] %s", r.Method, r.URL.Path))
}

func (s server) render(w http.ResponseWriter, r *http.Request) {
	session, err := s.session(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Render()

	if session.GameOver {
		w.Header().Set("HX-Reswap", "outerHTML")
	}

	err = s.t.ExecuteTemplate(w, "game", session)
	if err != nil {
		s.log.Error(err.Error())
	}

	s.log.Info(fmt.Sprintf("[%s] %s", r.Method, r.URL.Path))
}

func (s server) direction(w http.ResponseWriter, r *http.Request) {
	session, err := s.session(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := r.URL.Query().Get("keyTrigger")
	s.log.Info(fmt.Sprintf("[%s] %s", r.Method, r.URL.Path), "newKey", key)

	switch key {
	case arrowUpKey:
		session.SetDirection(game.Up)
	case arrowLeftKey:
		session.SetDirection(game.Left)
	case arrowDownKey:
		session.SetDirection(game.Down)
	case arrowRightKey:
		session.SetDirection(game.Right)
	}
}
