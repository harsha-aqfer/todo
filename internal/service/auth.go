package service

import (
	"github.com/gorilla/mux"
	"github.com/harsha-aqfer/todo/internal/db"
	"github.com/harsha-aqfer/todo/internal/util"
	"net/http"
	"strconv"
	"strings"
)

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, map[string]string{"error": "permission denied"})
}

func withApiKeyAuth(hf http.HandlerFunc, c *Config, u db.UserDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			permissionDenied(w)
			return
		} else {
			apiKey := authHeader[1]

			userEmail := mux.Vars(r)["user"]
			if userEmail == "" {
				permissionDenied(w)
				return
			}

			id, err := u.GetUserId(userEmail)
			if err != nil {
				permissionDenied(w)
				return
			}

			e := util.ComputeHmac256(strconv.Itoa(int(id)), c.SecretKey)
			if e != apiKey {
				permissionDenied(w)
				return
			}
			hf(w, r)
		}
	}
}
