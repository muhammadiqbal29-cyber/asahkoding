### ⚙️ DEEP-DIVE DEVOPS LIFECYCLE CHECKLIST (Enterprise Edition)

**1. Plan (Perencanaan)**
 * [x] Mengumpulkan Kebutuhan Bisnis (*Requirements Gathering*)
   * [x] Menentukan target performa (misal: *response time* aplikasi < 200ms)
   * [x] Menentukan estimasi jumlah pengguna aktif harian (*Daily Active Users*)
 * [x] Menyusun *Backlog* Produk & Prioritas Fitur
   * [x] Memisahkan fitur utama (*Must-Have*) dan fitur pendukung (*Nice-to-Have*)
 * [x] Membuat *User Story* dan Kriteria Penerimaan (*Acceptance Criteria*)
   * [x] Menulis skenario sukses dan skenario gagal (*edge cases*) untuk setiap fitur
 * [x] Menentukan *Timeline*, Estimasi Rilis, dan Alokasi Tim
   * [x] Menentukan siapa yang bertanggung jawab atas *code*, *infra*, dan *testing*
 * [x] **[NEW] Keandalan & Keamanan**
   * [x] Mendefinisikan SLA/SLO/SLI (Target Uptime, misalnya 99.9%)
   * [x] Membuat *Disaster Recovery Plan* (Skenario pemulihan jika server utama mati total)
   * [x] Melakukan *Threat Modeling* (Memetakan potensi celah keamanan di arsitektur aplikasi)

**2. Code (Pengodean)**
 * [x] **Pemilihan Teknologi & Bahasa Pemrograman**
   * [x] Menentukan bahasa/kerangka kerja untuk *Backend* / API Server (**Golang** dengan struktur folder ramah AI Context)
   * [x] Menentukan bahasa/kerangka kerja untuk *Frontend* (Antarmuka) (**Next.js**)
   * [x] Menentukan *Database* dan *Message Broker* (**MySQL** & Redis)
 * [x] **Perancangan Database & Skema**
   * [x] Membuat Entity Relationship Diagram (ERD) untuk tabel-tabel utama
 * [x] Inisialisasi Repositori Kode (Git Setup)
   * [x] Membuat file `.gitignore` (memastikan file rahasia tidak ikut terunggah)
 * [x] Konfigurasi Otoritas Keamanan & Hak Akses *(Dilewati - Proyek Solo)*
   * [x] Mengaktifkan *Branch Protection* di GitHub/GitLab (cabang `main`/`production` tidak bisa langsung di-*push* tanpa persetujuan)
   * [x] Mengunci akses token repositori hanya untuk tim yang berkepentingan
 * [x] **Persiapan Infrastruktur Development (Docker)**
   * [x] Membuat `docker-compose.yml` khusus untuk menjalankan MySQL & Redis lokal
 * [x] **Fase 1: Penulisan Kode Program Backend (Golang)**
   * [x] Membangun struktur folder Golang berbasis Domain (AI-Friendly)
   * [x] Menghubungkan Golang ke kontainer MySQL & Redis menggunakan `.env`
   * [x] Membangun *Execution Engine*: Golang mengendalikan sistem Docker untuk menjalankan kode peserta ujian
   * [x] Membuat fungsi koneksi *failover* (jika MySQL/Redis mati, API merespons dengan baik)
   * [x] **Keamanan Lapis Aplikasi (Autentikasi & Otorisasi)**
     * [x] Membuat sistem Login/Register (JWT Token)
     * [x] Membuat Middleware Otorisasi (Membedakan hak akses Admin vs Peserta Ujian)
   * [x] Membangun Rute API CRUD (Problem, Submission) & Skema Database
 * [x] **Fase 2: Penulisan Kode Program Frontend (Next.js)**
   * [x] Merancang antarmuka UI/UX (Daftar Soal, Editor Kode, Hasil Ujian)
   * [x] Menghubungkan Frontend dengan API Backend Golang
 * [x] Konfigurasi Deployment Lokal (Pengujian)
   * [x] Menjalankan Frontend dan Backend secara *Native* (tanpa dibungkus Docker) agar proses ngoding cepat (*hot-reload*)
 * [x] **[NEW] Standardisasi, Kualitas, & Keamanan Kode (DevSecOps)**
   * [x] Menerapkan *Linting* & *Formatting* otomatis (ESLint & Prettier)
   * [x] **[SECURITY]** Menambahkan pemindaian SAST (*Static Application Security Testing*) seperti `gosec` untuk mendeteksi celah keamanan (SQL Injection, Hardcoded Passwords) langsung di *source code*.
   * [x] Menambahkan Git Hooks (Husky) untuk mencegah *commit* jika kode tidak lulus *Linter* dan SAST.
   * [x] Menggunakan *Conventional Commits* (misal: `feat: auth`, `fix: login bug`) agar *history* rapi *(Dilewati - Proyek Solo)*
 * [x] Penulisan Pengujian Unit (*Unit Test*)
   * [x] Menulis tes untuk fungsi logika/matematika utama aplikasi
   * [x] Membuat *mocking* untuk *query* database agar tes berjalan cepat tanpa menyentuh DB asli
 * [x] Tinjauan Kode oleh Sesama Developer (*Peer Code Review*) *(Dilewati - Proyek Solo)*
   * [x] Minimal ada 1 atau 2 developer lain yang menyetujui (*approve*) *Pull Request* sebelum digabung
 * [x] Penggabungan Kode ke Cabang Utama (*Merge/Pull Request*) *(Dilewati - Proyek Solo)*

**3. Build (Pembuatan)**
 * [x] Otomatisasi pemicu pembuatan *build* menggunakan Jenkins (CI)
   * [x] Mengonfigurasi *Webhook* dari Git ke Jenkins (setiap ada *merge*, Jenkins otomatis jalan)
 * [x] Membuat Dockerfile
   * [x] Menggunakan *base image* yang ringan (misal: `node:alpine` atau `node:slim`) untuk menghemat ruang penyimpanan server
   * [x] Memasukkan perintah instalasi PM2 secara global (`RUN npm install pm2 -g`) *(Diganti dengan eksekusi langsung via kontainer)*
   * [x] **[NEW]** Menerapkan *Multi-stage Build* (Memisahkan proses instalasi `devDependencies` dan *runtime* agar *image* lebih kecil & aman)
 * [x] Menyusun *container image* menggunakan Docker
   * [x] Menjalankan perintah `docker build` lewat Jenkins
   * [x] **[NEW]** Memanfaatkan *Docker Cache Layer* untuk mempercepat proses *build* di pipeline
 * [ ] Melakukan kompilasi dan pembungkusan kode menjadi *Artifact*
   * [ ] Memastikan tidak ada *error* sintaksis saat proses *build* berjalan
   * [ ] **[SECURITY]** Melakukan *Container Security Scanning* (menggunakan Trivy/Clair) pada Docker Image yang baru dibuat untuk mendeteksi *vulnerability* level OS sebelum di-*deploy*.

**4. Test (Pengujian)**
 * [ ] Eksekusi *Unit Test* secara otomatis via pipeline Jenkins
   * [ ] Jika ada 1 saja *test* yang gagal (*fail*), Jenkins harus otomatis menghentikan proses (*abort pipeline*)
   * [ ] **[NEW]** Memeriksa *Code Coverage Threshold* (Gagal jika cakupan tes di bawah target, misal wajib 80%)
 * [ ] Menjalankan pengujian integrasi (*Integration Test*)
   * [ ] Menguji apakah kontainer aplikasi Node.js benar-benar bisa membaca dan menulis data ke kontainer Redis *test*
 * [ ] **[NEW] Pengujian Lanjutan (Advanced Testing)**
   * [ ] Menjalankan *Load / Stress Testing* (K6 / JMeter) untuk menguji ketahanan aplikasi saat beban tinggi
   * [ ] Menjalankan *End-to-End (E2E) Testing* secara otomatis pada antarmuka pengguna / alur sistem utuh
 * [ ] Pemindaian keamanan kode dan kualitas kode (*Security & Code Scanning*)
   * [ ] Memeriksa apakah ada pustaka *dependencies* yang memiliki celah keamanan (*vulnerability*) lewat `npm audit` / `govulncheck` atau integrasi SonarQube.
   * [ ] **[SECURITY]** Menjalankan pengujian DAST (*Dynamic Application Security Testing*) menggunakan OWASP ZAP untuk mensimulasikan serangan *hacker* (XSS, CSRF, Brute Force) ke API yang sedang menyala.

**5. Release (Rilis)**
 * [ ] Pemberian versi pada Docker Image (*Image Tagging*)
   * [ ] Melarang penggunaan tag `latest` di produksi. Harus menggunakan versi spesifik (misal: `:1.0.0` atau menggunakan *hash commit* Git)
   * [ ] **[NEW]** Menerapkan *Semantic Versioning* (SemVer) untuk standarisasi penomoran rilis (v1.0.0 ke v1.0.1, dst)
 * [ ] Mengunggah Docker Image aplikasi ke Docker Registry/Hub
   * [ ] Memastikan jalur pengunggahan (*push*) menggunakan koneksi yang aman/terenkripsi
 * [ ] **[NEW] Dokumentasi Rilis**
   * [ ] Menghasilkan *Automated Changelog* / *Release Notes* secara otomatis dari riwayat pesan *commit*
 * [ ] Menentukan rilis versi aplikasi yang stabil

**6. Deploy (Penerapan)**
 * [ ] **[NEW] Infrastruktur Otomatis & Strategi Rilis**
   * [ ] Menerapkan *Infrastructure as Code* (IaC) seperti Terraform / Ansible untuk *provisioning* server
   * [ ] Menyiapkan otomatisasi *Database Migrations* dalam alur penyebaran
   * [ ] Memilih Strategi *Zero-Downtime Deployment* (*Blue-Green Deployment*, *Rolling Updates*, atau *Canary Release*)
 * [ ] **[SECURITY] Infrastruktur Keamanan Jaringan & Rahasia**
   * [ ] Memasang WAF (*Web Application Firewall*) & *Rate Limiting* (contoh: Cloudflare/NGINX) untuk menangkis serangan DDoS dan *Brute Force* API *Login*.
   * [ ] Menggunakan *Secrets Management* (seperti HashiCorp Vault) untuk menyimpan *password* MySQL/Redis di produksi, **bukan** menggunakan file `.env` biasa.
 * [ ] Otomatisasi penyebaran (*Deployment*) menggunakan pipeline CD (Jenkins)
   * [ ] Mengamankan kunci SSH server produksi di dalam kredensial rahasia Jenkins (*Jenkins Credentials Store*)
 * [ ] Menyebarkan dan mengonfigurasi server/container Redis di lingkungan produksi
   * [ ] **Wajib:** Mengaktifkan password pada Redis produksi (jangan biarkan *default* tanpa *password*)
   * [ ] Membatasi *port* Redis (6379) agar hanya bisa diakses oleh internal container aplikasi Node.js, tidak terbuka untuk publik luar
 * [ ] Menjalankan container Docker berisi aplikasi Node.js di server produksi
   * [ ] Mengonfigurasi parameter *restart policy* di Docker (misal: `--restart unless-stopped`)
 * [ ] Memastikan perintah utama di dalam container dijalankan oleh PM2 (`pm2-runtime`)

**7. Operate (Pengoperasian)**
 * [ ] Mengelola jalannya kontainerisasi aplikasi dan database via Docker
   * [ ] Mengatur alokasi batas maksimum RAM dan CPU untuk masing-masing container agar tidak saling berebut *resource* server
 * [ ] Membiarkan PM2 menjaga proses Node.js tetap hidup
   * [ ] Memastikan PM2 berjalan dalam *Cluster Mode* jika server memiliki CPU lebih dari 1 *core*
 * [ ] Optimasi performa memori, replikasi, dan *cluster* Redis
   * [ ] Mengatur kebijakan *eviction* Redis (misal: `volatile-lru`) agar memori Redis tidak penuh dan menyebabkan server hang
 * [ ] **[NEW] Backup & Ketahanan Data**
   * [ ] Menerapkan *Backup & Restore Strategy* terjadwal untuk database dan Redis (Misal kirim *snapshot* ke Cloud Storage)
 * [ ] **[NEW] Pemusatan Log (Log Management)**
   * [ ] Mengirim seluruh log Node.js, PM2, dan Docker ke sistem terpusat (contoh: Elasticsearch/Kibana, atau Grafana Loki)
 * [ ] Manajemen skalabilitas infrastruktur (Auto-scaling) jika terjadi lonjakan *traffic* secara tiba-tiba

**8. Monitor (Pemantauan)**
 * [ ] Memantau kesehatan internal Node.js menggunakan PM2
   * [ ] Memeriksa jumlah *restart* yang dilakukan PM2 (jika angka *restart* tinggi, berarti ada *bug* fatal di dalam kode)
 * [ ] Memantau kesehatan, penggunaan *resource*, dan status container Docker
   * [ ] Mengatur peringatan (*alert*) otomatis (misal via Email/Slack) jika penggunaan CPU server menyentuh angka 85%
 * [ ] Memantau *resource*, memori, *hit rate*, dan latensi pada Redis
 * [ ] Memantau performa serta log kegagalan pada pipeline otomatis Jenkins
 * [ ] **[NEW] Metrik & Observabilitas Lanjutan**
   * [ ] Menerapkan *Application Performance Monitoring* (APM) dengan Dasbor (Prometheus & Grafana, New Relic, atau Datadog)
   * [ ] Memasang *Distributed Tracing* (seperti Jaeger/Zipkin) untuk melacak *bottleneck* performa di tiap layanan
   * [ ] Mengaktifkan *Uptime Monitoring* dari pihak eksternal (misal: UptimeRobot) untuk memverifikasi aplikasi online dari internet
