package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gegen07/cartola-university/application"
	"github.com/gegen07/cartola-university/domain/entity/scout"
	"github.com/gegen07/cartola-university/interfaces/errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ScoutHandler struct {
	ScoutApplication application.ScoutApplicationInterface
}

func NewScoutHandler(scoutApplication application.ScoutApplicationInterface) *ScoutHandler {
	return &ScoutHandler{
		ScoutApplication: scoutApplication,
	}
}

func (handler *ScoutHandler) GetAllScouts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	w.Header().Set("Content-Type", "application/json")

	key := r.FormValue("page")
	page, err := strconv.Atoi(key)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid Parameter")
	}

	var scouts scout.Scouts
	scouts, err = handler.ScoutApplication.GetAll(r.Context(), page)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scouts.PublicScouts())

	return nil
}

func (handler *ScoutHandler) GetScoutByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid JSON")
	}

	scout, err := handler.ScoutApplication.GetByID(r.Context(), id)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(scout.PublicScout())

	return nil
}

func (handler *ScoutHandler) Insert(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	var scout scout.RequestScout
	err := json.NewDecoder(r.Body).Decode(&scout)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid JSON")
	}

	s, err := handler.ScoutApplication.Insert(r.Context(), &scout)

	if err != nil {
		return fmt.Errorf("DB Error: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s.PublicScout())
	return nil
}

func (handler *ScoutHandler) Update(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPut {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
	var scout scout.Scout
	var err error

	params := mux.Vars(r)
	scout.ID, err = strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid JSON")
	}

	err = json.NewDecoder(r.Body).Decode(&scout)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid JSON")
	}

	s, err := handler.ScoutApplication.Update(r.Context(), &scout)

	if err != nil {
		return fmt.Errorf("DB Error: %v", err)
	}

	json.NewEncoder(w).Encode(s.PublicScout())
	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler *ScoutHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid JSON")
	}

	err = handler.ScoutApplication.Delete(r.Context(), id)

	if err != nil {
		return fmt.Errorf("DB Error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
