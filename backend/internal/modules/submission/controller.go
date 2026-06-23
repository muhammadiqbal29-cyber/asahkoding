package submission

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/muhammad-iqbal/leetcode-backend/internal/config"
	"github.com/muhammad-iqbal/leetcode-backend/internal/middleware"
	"github.com/muhammad-iqbal/leetcode-backend/internal/modules/problem"
)

type SubmitRequest struct {
	ProblemID uint   `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

// SubmitCode adalah otak aplikasi: menerima kode, mengujinya, lalu memberikan nilai
func SubmitCode(w http.ResponseWriter, r *http.Request) {
	// 1. Validasi Identitas (Pastikan yang akses sudah Login)
	claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Catat riwayat awal ke database sebagai "Pending"
	sub := Submission{
		UserID:    claims.UserID,
		ProblemID: req.ProblemID,
		Language:  req.Language,
		Code:      req.Code,
		Status:    "Pending",
	}
	config.DB.Create(&sub)

	// 3. Ambil Test Case rahasia dari Soal tersebut
	var p problem.Problem
	if err := config.DB.Preload("TestCases").First(&p, req.ProblemID).Error; err != nil {
		sub.Status = "Error: Problem not found"
		config.DB.Save(&sub)
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	// 4. Pilih Mesin Juri yang sesuai bahasa pemrogramannya
	var runner CodeRunner
	if req.Language == "go" {
		runner = &GoRunner{}
	} else {
		sub.Status = "Error: Unsupported language"
		config.DB.Save(&sub)
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// === FASE 1: KOMPILASI KODE ===
	compileDir, err := runner.Compile(req.Code)
	if err != nil {
		sub.Status = "Compile Error"
		config.DB.Save(&sub)
		
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"submission_id": sub.ID,
			"status":        "Compile Error",
			"error":         err.Error(),
		})
		return
	}
	// Pastikan folder kompilasi dihapus setelah penilaian selesai
	defer runner.Cleanup(compileDir)

	// === FASE 2: UJIAN TEST CASE ===
	allPassed := true
	var maxTime int64 = 0
	var failureReason string = ""

	for i, tc := range p.TestCases {
		// Lempar Input soal ke dalam Stdin kontainer Docker
		tcResult := runner.RunTestCase(compileDir, tc.Input, p.TimeLimitMs)
		
		if tcResult.TimeTakenMs > maxTime {
			maxTime = tcResult.TimeTakenMs
		}

		if tcResult.Error == "Time Limit Exceeded" {
			sub.Status = "Time Limit Exceeded"
			allPassed = false
			break
		}

		if tcResult.Error != "" {
			sub.Status = "Runtime Error"
			failureReason = tcResult.Error
			allPassed = false
			break
		}

		// Penjurian: Apakah hasil cetak terminal sama persis dengan kunci jawaban?
		if strings.TrimSpace(tcResult.Output) != strings.TrimSpace(tc.ExpectedOutput) {
			sub.Status = "Wrong Answer"
			failureReason = "Test Case " + string(rune(i+1)) + " Failed"
			allPassed = false
			break
		}
	}

	if allPassed {
		sub.Status = "Accepted"
	}

	// 5. Simpan nilai rapot akhir ke MySQL
	sub.ExecutionTime = int(maxTime)
	config.DB.Save(&sub)

	// 6. Kembalikan hasil ujian ke User
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"submission_id":  sub.ID,
		"status":         sub.Status,
		"execution_time": sub.ExecutionTime,
		"error_detail":   failureReason,
	})
}
