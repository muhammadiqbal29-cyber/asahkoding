package submission

type TestCaseResult struct {
	Output      string
	Error       string
	TimeTakenMs int64
}

// CodeRunner adalah mesin Juri Otomatis (Dua Fase)
type CodeRunner interface {
	// Fase 1: Menerjemahkan kode menjadi biner (mengembalikan path folder rahasia)
	Compile(code string) (string, error)
	
	// Fase 2: Menjalankan biner dengan Input spesifik (via stdin) di dalam Docker
	RunTestCase(compileDir string, input string, timeLimitMs int) TestCaseResult
	
	// Fase 3: Membersihkan folder setelah semua test case selesai
	Cleanup(compileDir string)
}
