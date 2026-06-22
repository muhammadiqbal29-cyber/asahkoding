# Product Requirements Document (PRD)
## Project: LeetCode Clone (Lokal)

### 1. Deskripsi Proyek
Sebuah platform evaluasi dan pembelajaran kode interaktif yang menyerupai LeetCode. Platform ini memungkinkan pengguna untuk memecahkan masalah pemrograman dengan menulis, mengirimkan (*submit*), dan mengeksekusi kode secara langsung di peramban web (*browser*). Aplikasi ini di-host secara lokal dari sebuah laptop pribadi dan diekspos ke publik internet.

### 2. Tujuan Bisnis & Target Pengguna
*   **Target Pasar:** Pengguna internet publik, dikhususkan untuk pengguna di wilayah Indonesia.
*   **Ketersediaan (SLA/Uptime):** Aplikasi beroperasi secara *ad-hoc* (0% SLA terjamin). Aplikasi hanya akan *online* dan dapat diakses ketika server/laptop utama (Host) dinyalakan oleh pemilik.

### 3. Spesifikasi Fungsional & Performa
*   **Dukungan Bahasa Pemrograman (Fase 1):**
    *   Go (Compiled)
    *   PHP (Interpreted)
    *   Python (Interpreted)
    *   JavaScript / Node.js (Interpreted)
*   **Kinerja Eksekusi Kode (*Response Time*):**
    *   Sistem ditargetkan mengembalikan hasil *output* (*Success/Fail/Error*) dalam waktu **di bawah 3 detik** sejak pengguna menekan tombol *Submit*.

### 4. Spesifikasi Non-Fungsional (Infrastruktur & Keamanan)
*   **Beban Pengguna Berbarengan (*Concurrency*):**
    *   Sistem dirancang untuk sanggup menerima beban hingga batas maksimum **500 eksekusi/submit bersamaan**. Sistem harus memprioritaskan antrean agar laptop/host tidak mengalami gagal sistem (*crash*) akibat kehabisan memori atau CPU.
*   **Keamanan Eksekusi (Isolasi):**
    *   Kode pengguna yang tidak tepercaya (*untrusted code*) wajib dijalankan di dalam lingkungan terisolasi menggunakan **Docker Container**.
    *   Kontainer ini bersifat sementara (*ephemeral*), tidak memiliki akses ke jaringan lokal/internet tanpa izin, dan akan langsung dihancurkan setelah eksekusi selesai untuk melindungi *file system* laptop.

### 5. Arsitektur Sistem & Tech Stack (Final)
Berdasarkan batasan perangkat keras (laptop lokal) dan keputusan pengembangan, ini adalah alat dan arsitektur yang digunakan:
1.  **Public Exposure (Tunneling):** **Cloudflare Tunnels** (atau Ngrok). Meneruskan trafik publik internet langsung ke *localhost* dengan aman (SSL/TLS) tanpa perlu *Port Forwarding*.
2.  **Frontend (UI/Antarmuka):** **Next.js**. Membangun antarmuka interaktif termasuk integrasi *code editor* langsung di peramban pengguna.
3.  **Backend (API Server):** **Golang**. Digunakan untuk performa tingkat tinggi (CPU/RAM yang ringan lewat *Goroutines*) saat menerima beban maksimal, serta kecepatan interaksi dengan API Docker. Menggunakan **struktur folder berbasis Domain/Feature (AI-Context Friendly)** untuk mempermudah navigasi.
4.  **Database Utama:** **MySQL**. Menyimpan data terstruktur yang memiliki relasi kuat seperti data *User*, *Problem* (Bank Soal), dan riwayat *Submission*.
5.  **Sistem Antrean (*Message Broker*):** **Redis**. Berfungsi menampung masuknya 500 *request* seketika. Worker Golang akan memproses *request* dari antrean Redis secara aman bagi CPU laptop.
6.  **Mesin Eksekusi (*Execution Engine*):** **Docker**. Dikelola oleh backend Golang untuk membuat dan membuang kontainer secara dinamis. Untuk target respons di bawah 3 detik, dapat diterapkan metode *Pre-warmed Containers* (menyiapkan kontainer kosong yang *standby*).
