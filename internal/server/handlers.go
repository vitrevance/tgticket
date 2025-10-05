package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vitrevance/tgticket/internal/ticket"
)

func (s *Server) adminHandler(w http.ResponseWriter, r *http.Request) {
	activeTickets := s.tstore.ActiveTickets()

	data := struct {
		Tickets       []*ticket.Ticket
		Now           time.Time
		PublicAddress string
	}{
		Tickets:       activeTickets,
		Now:           time.Now(),
		PublicAddress: s.cfg.PublicAddr,
	}

	err := s.tmpl.ExecuteTemplate(w, "admin.html", data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

func (s *Server) newTicketHandler(w http.ResponseWriter, r *http.Request) {
	token, err := s.generateToken(16)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	t := &ticket.Ticket{
		Token:    token,
		ExpireAt: time.Now().Add(6 * time.Hour),
	}
	s.tstore.Add(t)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) prolongTicketHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	t, ok := s.tstore.Get(token)
	if !ok {
		http.Error(w, "ticket not found", http.StatusNotFound)
		return
	}
	t.ExpireAt = t.ExpireAt.Add(6 * time.Hour)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) revokeTicketHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	s.tstore.Delete(token)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) controlHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[len("/control/"):]

	t, ok := s.tstore.Get(token)
	if !ok || t.ExpireAt.Before(time.Now()) {
		// Ticket expired - show expired ticket page
		if s.tmpl == nil {
			if err := s.parseTemplates(); err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := s.tmpl.ExecuteTemplate(w, "expired.html", nil)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.renderControlPage(w)
	case http.MethodPost:
		s.handleControlPost(w, r, token)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) renderControlPage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.tmpl.ExecuteTemplate(w, "control.html", nil)
}

func (s *Server) handleControlPost(w http.ResponseWriter, r *http.Request, token string) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
		return
	}
	num := r.PostForm.Get("gateNumber")
	if num == "" {
		http.Error(w, "Номер шлагбаума обязателен", http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("Нажата кнопка открытия шлагбаума\nВремя: %s\nНомер шлагбаума: %s\nID доступа: %s",
		time.Now().Format("2006-01-02 15:04:05"),
		num,
		token,
	)

	go func() {
		if err := s.sendTelegramMessage(msg); err != nil {
			log.Printf("Ошибка отправки сообщения Telegram: %v", err)
		}
	}()

	// Confirmation page inline
	const confPage = `
<!DOCTYPE html>
<html lang="ru">
<head><meta charset="UTF-8" /><title>Открытие начато</title></head>
<body>
<h1>Открытие шлагбаума</h1>
<p>Сигнал отправлен. Ожидайте 10-30 секунд.</p>
<p><a href="">Вернуться</a></p>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, confPage)
}
