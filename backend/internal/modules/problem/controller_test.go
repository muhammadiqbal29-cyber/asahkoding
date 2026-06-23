package problem

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
	"github.com/muhammad-iqbal/leetcode-backend/internal/testutils"
)

func TestGetProblems(t *testing.T) {
	testutils.SetupMockDB()
	config.DB.AutoMigrate(&Problem{}, &TestCase{})

	// 1. Seed Data Soal ke DB RAM
	config.DB.Create(&Problem{Title: "Soal Mudah", Description: "Mudah", Difficulty: "Easy"})
	config.DB.Create(&Problem{Title: "Soal Sulit", Description: "Sulit", Difficulty: "Hard"})

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProblems)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetProblems salah kode status: didapat %v seharusnya %v", status, http.StatusOK)
	}

	var problems []Problem
	json.NewDecoder(rr.Body).Decode(&problems)
	if len(problems) != 2 {
		t.Errorf("GetProblems gagal mereturn 2 soal, malah mereturn %v", len(problems))
	}
}

func TestGetProblemByID(t *testing.T) {
	testutils.SetupMockDB()
	config.DB.AutoMigrate(&Problem{}, &TestCase{})

	// 1. Seed Data Soal dengan Test Cases
	p := Problem{
		Title: "Soal Test Case",
		TestCases: []TestCase{
			{Input: "1", ExpectedOutput: "1", IsHidden: false},
			{Input: "2", ExpectedOutput: "2", IsHidden: true}, // Test case rahasia (TIDAK BOLEH BOCOR)
		},
	}
	config.DB.Create(&p)

	// Buat simulasi request API dengan Chi URL Param (id=1)
	req, _ := http.NewRequest("GET", "/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProblemByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetProblemByID salah kode status: didapat %v seharusnya %v", status, http.StatusOK)
	}

	var res Problem
	json.NewDecoder(rr.Body).Decode(&res)
	
	// Pastikan hanya test case publik yang keluar
	if len(res.TestCases) != 1 {
		t.Errorf("KEBOCORAN SOAL: Diharapkan hanya 1 test case publik yang keluar, tapi muncul %d", len(res.TestCases))
	}
	if res.TestCases[0].Input != "1" {
		t.Errorf("Test case yang keluar salah")
	}
}
