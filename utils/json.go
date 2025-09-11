package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"log"
	"net/http"
	"strings"
)

type JsonResponse[T any] struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    T      `json:"data"`
}

func RespondWithJson(w http.ResponseWriter, code int, content interface{}) {

	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(content)

	if err != nil {
		log.Printf("Unable to marshal json %v", content)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message string, data interface{}) {

	w.Header().Add("Content-Type", "application/json")

	if data == nil {
		data = struct{}{}
	}

	content := JsonResponse[interface{}]{
		message,
		false,
		data,
	}

	marshalledData, err := json.Marshal(content)

	if err != nil {
		log.Printf("Unable to marshal json %v", content)
		w.Write([]byte(`{"message": "Internal Server Error", "success": false, "data": {}}`))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(marshalledData)
}

func RespondWithSuccess(w http.ResponseWriter, code int, message string, data interface{}) {

	w.Header().Add("Content-Type", "application/json")

	if data == nil {
		data = struct{}{}
	}

	content := JsonResponse[interface{}]{
		message,
		true,
		data,
	}

	marshalledData, err := json.Marshal(content)

	if err != nil {
		log.Printf("Unable to marshal json %v", content)
		w.Write([]byte(`{"message": "Internal Server Error", "success": false, "data": {}}`))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(marshalledData)
}

func JsonValidate[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var payload T

	// Decode JSON body
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		RespondWithError(w, 400, "Invalid JSON", nil)
		return payload, err
	}

	validate := validator.New()
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")

	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		RespondWithError(w, 400, "Invalid JSON", nil)
		return payload, err
	}

	err = validate.Struct(payload)
	if err != nil {
		// Validation failed, handle the error
		messages := make([]string, 0)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				messages = append(messages, fe.Translate(trans))
			}
		}
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", strings.Join(messages, ", ")), nil)
		return payload, err
	}

	return payload, nil
}
