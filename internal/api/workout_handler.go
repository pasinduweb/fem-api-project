package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pasinduweb/fem-api-project/internal/store"
	"github.com/pasinduweb/fem-api-project/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutId)

	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: decodeCreateWorkout: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request sent"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)

	if err != nil {
		wh.logger.Printf("ERROR: createWorkout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "failed to create workout"})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envelop{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout update id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)

	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// Assume we are able to find an existing workout
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		wh.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request payload"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)

	if err != nil {
		wh.logger.Printf("ERROR: updatingWorkout: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")

	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)

	if err == sql.ErrNoRows {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Error deleting workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
