package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type Helper struct{}

// Helper function to write an error response
func (h *Helper) WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorResponse{
		Success: false,
		Message: message,
		Data:    err,
	})
}

// Helper function to write a success response
func (h *Helper) WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// Helper function that turns a string into an integer
func (h *Helper) StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("failed to convert string to int")
	}

	return i, nil
}

// Helper function that turns a string into an primitive.ObjectID
func (h *Helper) StringToPrimitiveObjectID(s string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.ObjectID{}, errors.New("failed to convert string to primitive.ObjectID")
	}

	return id, nil
}
