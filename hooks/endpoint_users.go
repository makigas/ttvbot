package hooks

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	"gopkg.makigas.es/redismemo"
	"gopkg.makigas.es/ttvbot/httpd/server"
	"gopkg.makigas.es/ttvbot/integrations/ttv/helixapi"
)

func EndpointUsers(
	httpd *server.HttpServer,
	hapi *helixapi.HelixApi,
	rdb *redis.Client,
) {
	memo := redismemo.RedisMemo(rdb)
	httpd.AddHandler(func(router *chi.Mux) {
		router.Get("/helix/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if name := r.URL.Query().Get("name"); name != "" {
				// We have a name, we want the ID.
				key := "helix:name_to_id:" + name

				computeId := computeUserId{hapi: hapi, name: name}
				val, err := memo(r.Context(), key, computeId.compute, 24*time.Hour)
				if err != nil {
					http.Error(w, "Cannot fetch user ID: "+err.Error(), http.StatusBadGateway)
					return
				}
				if computeId.err != nil {
					http.Error(w, "Cannot fetch user ID: "+err.Error(), http.StatusBadGateway)
					return
				}

				// Here is the user ID.
				var response struct {
					UserName string
					UserId   string
				}
				response.UserName = name
				response.UserId = val
				data, err := json.Marshal(&response)
				if err != nil {
					http.Error(w, "Serialization error: "+err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Add("Content-Type", "application/json")
				w.Write(data)
				return
			}

			http.Error(w, "I do not know what do you want from me", http.StatusBadGateway)
		}))
	})
}

type computeUserId struct {
	hapi *helixapi.HelixApi
	name string
	err  error
}

func (comp *computeUserId) compute() string {
	client, err := comp.hapi.NewAppClient()
	if err != nil {
		comp.err = err
		return ""
	}

	resp, err := client.GetUsers(&helix.UsersParams{
		Logins: []string{comp.name},
	})
	if err != nil {
		comp.err = err
		return ""
	}
	if len(resp.Data.Users) != 1 {
		comp.err = errors.New("expected data from an user to have been returned")
		return ""
	}
	return resp.Data.Users[0].ID
}
