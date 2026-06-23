package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
	"github.com/muhammad-iqbal/leetcode-backend/internal/testutils"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterAndLogin(t *testing.T) {
	testutils.SetupMockDB()
	config.DB.AutoMigrate(&User{})

	// 1. UJI API REGISTER
	reqBody, _ := json.Marshal(RegisterRequest{
		Username: "fauzan",
		Email:    "fauzan@leetcode.com",
		Password: "rahasia123",
	})
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Register salah kode status: didapat %v seharusnya %v", status, http.StatusCreated)
	}

	// Buktikan password di-hash di dalam DB Mock
	var u User
	config.DB.First(&u)
	if u.Username != "fauzan" {
		t.Errorf("User gagal disimpan ke DB")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("rahasia123")); err != nil {
		t.Errorf("Keamanan jebol: Password gagal di-hash dengan bcrypt")
	}

	// 2. UJI API LOGIN
	loginBody, _ := json.Marshal(LoginRequest{
		Email:    "fauzan@leetcode.com",
		Password: "rahasia123",
	})
	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	rrLogin := httptest.NewRecorder()
	handlerLogin := http.HandlerFunc(LoginHandler)
	handlerLogin.ServeHTTP(rrLogin, reqLogin)

	if status := rrLogin.Code; status != http.StatusOK {
		t.Errorf("Login salah kode status: didapat %v seharusnya %v", status, http.StatusOK)
	}

	var response map[string]string
	json.NewDecoder(rrLogin.Body).Decode(&response)
	if response["token"] == "" {
		t.Errorf("Sistem gagal menerbitkan Token JWT")
	}
}
