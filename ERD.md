# Entity Relationship Diagram (ERD)
## LeetCode Clone Lokal

```mermaid
erDiagram
    USERS {
        int id PK
        string username
        string email
        string password_hash
        string role "admin/user"
        datetime created_at
    }

    PROBLEMS {
        int id PK
        string title
        text description
        string difficulty "Easy/Medium/Hard"
        int time_limit_ms
        int memory_limit_kb
        datetime created_at
    }

    TEST_CASES {
        int id PK
        int problem_id FK
        text input_data
        text expected_output
        boolean is_hidden "Untuk test rahasia"
    }

    SUBMISSIONS {
        int id PK
        int user_id FK
        int problem_id FK
        string language "go, php, python, js"
        text source_code
        string status "Pending, Success, Fail, Error"
        int execution_time_ms
        int memory_used_kb
        datetime submitted_at
    }

    USERS ||--o{ SUBMISSIONS : "melakukan"
    PROBLEMS ||--o{ SUBMISSIONS : "menerima"
    PROBLEMS ||--|{ TEST_CASES : "memiliki"
```

### Penjelasan Tabel Utama (MySQL)
1. **USERS**: Menyimpan data autentikasi pengguna. Memiliki kolom `role` untuk membedakan **Admin** (pemilik yang bisa membuat/menghapus soal) dan **User** (peserta biasa).
2. **PROBLEMS**: Bank soal. Menyimpan detail soal, tingkat kesulitan, serta batasan waktu (*time limit*) dan memori yang diizinkan saat eksekusi di Docker.
3. **TEST_CASES**: Menyimpan *input* dan *expected output* (harapan hasil). Terdapat kolom `is_hidden` agar kita bisa membuat *test case* yang disembunyikan dari pengguna untuk mencegah mereka melakukan kecurangan (misal *hardcode* jawaban).
4. **SUBMISSIONS**: Tabel paling krusial. Menyimpan setiap kode yang disubmit, bahasa yang dipakai, status kelulusan *test case*, serta seberapa cepat kode tersebut berjalan.

*(Catatan: Redis tidak dimasukkan ke dalam ERD karena Redis bersifat penyimpanan sementara / cache / queue, bukan penyimpanan relasional permanen).*
