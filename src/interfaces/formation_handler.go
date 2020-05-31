package interfaces

import (
	"encoding/json"
	"fmt"
	"github.com/gegen07/cartola-university/application"
	"github.com/gegen07/cartola-university/domain/entity"
	"github.com/gegen07/cartola-university/interfaces/errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type FormationHandler struct {
	formationApplication application.FormationApplicationInterface
}

func NewFormationHandler(formationApp application.FormationApplicationInterface) *FormationHandler {
	return &FormationHandler{
		formationApplication: formationApp,
	}
}

func (handler *FormationHandler) GetAllFormations(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return errors.NewHTTPError(nil, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
	w.Header().Set("Content-Type", "application/json")

	key := r.FormValue("page")
	page, err := strconv.Atoi(key)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid Parameter")
	}

	var formations entity.Formations
	formations, err = handler.formationApplication.GetAll(r.Context(), page)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formations.PublicFormations())

	return nil
}

func (handler *FormationHandler) GetFormationById(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	formationId, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid JSON")
	}

	formation, err := handler.formationApplication.GetByID(r.Context(), formationId)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formation.PublicFormation())

	return nil
}

func (handler *FormationHandler) Insert(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var formation entity.Formation

	err := json.NewDecoder(r.Body).Decode(&formation)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid JSON")
	}

	f, err := handler.formationApplication.Insert(r.Context(), &formation)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(f.PublicFormation())

	return nil
}

func (handler *FormationHandler) Update(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var formation entity.Formation
	var err error

	params := mux.Vars(r)
	formation.ID, err = strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid JSON")
	}

	err = json.NewDecoder(r.Body).Decode(&formation)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid JSON")
	}

	f, err := handler.formationApplication.Update(r.Context(), &formation)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(f.PublicFormation())

	return nil
}

func (handler *FormationHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	formationId, err := strconv.ParseUint(params["id"], 10, 64)

	if err != nil {
		return errors.NewHTTPError(err, 400, "Invalid JSON")
	}

	err = handler.formationApplication.Delete(r.Context(), formationId)

	if err != nil {
		return fmt.Errorf("DB error: %v", err)
	}

	w.WriteHeader(http.StatusOK)

	return nil
}