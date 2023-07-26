package eventsub

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gopkg.makigas.es/ttvbot/httpd/server"
)

var Module = fx.Module("Eventsub", fx.Provide(NewEventSubManager))

type EventHandler func(payload map[string]interface{})

type EventSubManager struct {
	rdb       *redis.Client
	listeners map[string]EventHandler
}

func NewEventSubManager(rdb *redis.Client, httpd *server.HttpServer) *EventSubManager {
	manager := EventSubManager{rdb: rdb, listeners: make(map[string]EventHandler)}
	httpd.AddHandler(func(mux *chi.Mux) {
		mux.Route("/eventsub", func(r chi.Router) {
			r.Use(manager.requestValidatorMiddleware)
			r.Use(manager.preventDuplicationMiddleware)
			r.Post("/", manager.handleEvent)
		})
	})
	return &manager
}

func (es *EventSubManager) AddEventListener(kind string, handler EventHandler) {
	es.listeners[kind] = handler
}

func (es *EventSubManager) requestValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("Twitch-Eventsub-Subscription-Type")
		messageId := r.Header.Get("Twitch-Eventsub-Message-Id")
		timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
		signature := r.Header.Get("Twitch-Eventsub-Message-Signature")

		// Get the secret assigned to this event subscription type.
		secret, err := es.rdb.Get(context.Background(), "eventsub:secrets:"+eventType).Result()
		if err != nil && err != redis.Nil {
			http.Error(w, "Internal error while validating request", http.StatusBadGateway)
			return
		}

		// Read the message body, required to build a signature.
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Check the message signature.
		payload := messageId + timestamp + string(body)
		hasher := hmac.New(sha256.New, []byte(secret))
		hasher.Write([]byte(payload))
		selfsign := "sha256=" + hex.EncodeToString(hasher.Sum(nil))
		if signature != selfsign {
			http.Error(w, "Invalid HMAC signature.", http.StatusForbidden)
			return
		}

		// Signature is valid.
		next.ServeHTTP(w, r)
	})
}

func (es *EventSubManager) preventDuplicationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messageId := r.Header.Get("Twitch-Eventsub-Message-Id")
		messageKey := "eventsub:id:" + messageId

		// Check if we have already handled this ID.
		count, err := es.rdb.Exists(context.Background(), messageKey).Result()
		if err != nil {
			http.Error(w, "Internal error while validating request", http.StatusBadGateway)
			return
		}
		if count > 0 {
			http.Error(w, "Duplicated identifier: did I handle this?", http.StatusBadRequest)
			return
		}

		// Handle the notification quickly.
		next.ServeHTTP(w, r)

		// Put the ID in the storage to prevent double notifications.
		es.rdb.Set(context.Background(), messageKey, messageId, 0)
	})
}

func (es *EventSubManager) handleEvent(w http.ResponseWriter, r *http.Request) {
	kind := r.Header.Get("Twitch-Eventsub-Message-Type")
	switch kind {
	case "webhook_callback_verification":
		es.handleChallenge(w, r)
	case "revocation":
		es.handleRevocation(w, r)
	case "notification":
		es.handleNotification(w, r)
	default:
		http.Error(w, "Not handling this type: "+kind, http.StatusNotImplemented)
	}
}

func (es *EventSubManager) handleChallenge(w http.ResponseWriter, r *http.Request) {
	// Decode payload.
	var payload struct {
		Challenge string
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Cannot decode challenge payload", http.StatusBadRequest)
		return
	}

	// Ack request.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload.Challenge))
}

func (es *EventSubManager) handleNotification(w http.ResponseWriter, r *http.Request) {
	// Decode payload.
	var payload struct {
		Subscription struct {
			Type string `json:"type"`
		} `json:"subscription"`
		Event map[string]interface{} `json:"event"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Cannot decode challenge payload", http.StatusBadRequest)
		return
	}

	// Ack request.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("have a nice day"))

	// Call the appropiate event listener if configured.
	if listener, ok := es.listeners[payload.Subscription.Type]; ok {
		listener(payload.Event)
	}
}

func (es *EventSubManager) handleRevocation(w http.ResponseWriter, r *http.Request) {
	// TODO: alert somewhere?
	data, _ := io.ReadAll(r.Body)
	fmt.Println(string(data))

	// Ack request.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Acknowledged revocation"))
}
