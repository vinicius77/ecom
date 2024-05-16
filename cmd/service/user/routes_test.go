package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vinicius77/ecom/types"
)

func TestUserServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("Should fail if payload is invalid", func(t *testing.T) {

		payload := types.RegisterUserPayload{
			FirstName: "vinicius",
			LastName:  "bonifacio",
			Email:     "",
			Password:  "XXXXXXXX",
		}

		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		// response recorder
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

	})
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user with email %s not found", email)
}

func (m *mockUserStore) GetUserById(id string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil

}
