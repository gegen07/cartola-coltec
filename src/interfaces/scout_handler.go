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
	PositionApplication application.PositionApplicationInterface
}

func NewScoutHandler(scoutApplication application.ScoutApplicationInterface, positionApplication application.PositionApplicationInterface) *ScoutHandler {
	return &ScoutHandler{
		ScoutApplication: scoutApplication,
		PositionApplication: positionApplication,
	}
}

func (handler *ScoutHandler) GetAllScouts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	w.Header().Set("Content-Type", "application/json")

	var scouts scout.Scouts
	scouts, err := handler.ScoutApplication.GetAll()

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

	scout, err := handler.ScoutApplication.GetByID(id)

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

	var scout scout.Scout
	err := json.NewDecoder(r.Body).Decode(&scout)

	if err != nil {
		return errors.NewHTTPError(err, http.StatusBadRequest, "Invalid JSON")
	}

	position, err := handler.PositionApplication.GetById(scout.ID)
	err = handler.PositionApplication.AppendScoutAssociation(position, &scout)
	s, err := handler.ScoutApplication.Insert(&scout)

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

	s, err := handler.ScoutApplication.Update(&scout)

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

	err = handler.ScoutApplication.Delete(id)

	if err != nil {
		return fmt.Errorf("DB Error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
