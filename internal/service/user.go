package service

import (
	"encoding/json"
	"github.com/harsha-aqfer/todo/internal/util"
	"github.com/harsha-aqfer/todo/pkg"
	"net/http"
	"strconv"
)

func (s *Service) HandleUserRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	var ui pkg.UserInfo

	if err := json.NewDecoder(r.Body).Decode(&ui); err != nil {
		return err
	}

	if err := ui.Validate(); err != nil {
		return apiError{msg: err.Error(), status: http.StatusBadRequest}
	}

	id, err := s.db.User.CreateUser(&ui)
	if err != nil {
		return err
	}

	//construct api key
	ui.ApiKey = util.ComputeHmac256(strconv.Itoa(int(id)), s.c.SecretKey)
	return WriteJSON(w, http.StatusOK, ui)
}
