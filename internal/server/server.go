package server

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/vitrevance/tgticket/internal/config"
	"github.com/vitrevance/tgticket/internal/ticket"
	"github.com/vitrevance/tgticket/templates"
)

type Server struct {
	cfg    config.Config
	tstore *ticket.Store
	tmpl   *template.Template
}

func NewServer(cfg config.Config) *Server {
	return &Server{
		cfg:    cfg,
		tstore: ticket.NewStore(),
	}
}

func (s *Server) RegisterRoutes() {
	err := s.parseTemplates()
	if err != nil {
		log.Fatal("Failed to parse templates", err)
	}
	http.HandleFunc("/login", s.loginHandler)
	http.HandleFunc("/logout", s.logoutHandler)
	http.Handle("/admin/", s.authMiddleware(http.HandlerFunc(s.adminHandler)))
	http.Handle("/admin/ticket/new", s.authMiddleware(http.HandlerFunc(s.newTicketHandler)))
	http.Handle("/admin/ticket/prolong", s.authMiddleware(http.HandlerFunc(s.prolongTicketHandler)))
	http.Handle("/admin/ticket/revoke", s.authMiddleware(http.HandlerFunc(s.revokeTicketHandler)))
	http.HandleFunc("/control/", s.controlHandler)
}

func (s *Server) parseTemplates() error {
	var err error
	s.tmpl, err = template.New("").Funcs(template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
	}).ParseFS(templates.TemplateFiles, "*.html")
	return err
}

func (s *Server) generateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// (handlers detailed in next snippets)
