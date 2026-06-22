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
 * [x] Inisialisasi Repositori Kode (Git Setup)
   * [x] Membuat file `.gitignore` (memastikan `node_modules` dan file `.env` tidak ikut terunggah)
 * [ ] Konfigurasi Otoritas Keamanan & Hak Akses
   * [ ] Mengaktifkan *Branch Protection* di GitHub/GitLab (cabang `main`/`production` tidak bisa langsung di-*push* tanpa persetujuan)
   * [ ] Mengunci akses token repositori hanya untuk tim yang berkepentingan
 * [ ] Penulisan Kode Program Node.js (Fitur & Integrasi Redis)
   * [ ] Menggunakan *environment variables* (`process.env`) untuk koneksi database, bukan *hardcode*
   * [ ] Membuat fungsi koneksi *failover* (jika Redis mati, aplikasi tidak langsung *crash*)
 * [ ] Membuat file konfigurasi PM2 (`ecosystem.config.js`)
   * [ ] Mengatur `instances: "max"` atau `max_memory_restart` agar PM2 otomatis melakukan *restart* jika memori Node.js bocor (*memory leak*)
 * [ ] **[NEW] Standardisasi & Kualitas Kode**
   * [ ] Menerapkan *Linting* & *Formatting* otomatis (ESLint & Prettier)
   * [ ] Menambahkan Git Hooks (Husky) untuk memeriksa *linter* sebelum *commit* (*Pre-commit hook*)
   * [ ] Menggunakan *Conventional Commits* (misal: `feat: auth`, `fix: login bug`) agar *history* rapi
 * [ ] Penulisan Pengujian Unit (*Unit Test*)
   * [ ] Menulis tes untuk fungsi logika/matematika utama aplikasi
   * [ ] Membuat *mocking* untuk *query* database agar tes berjalan cepat tanpa menyentuh DB asli
 * [ ] Tinjauan Kode oleh Sesama Developer (*Peer Code Review*)
   * [ ] Minimal ada 1 atau 2 developer lain yang menyetujui (*approve*) *Pull Request* sebelum digabung
 * [ ] Penggabungan Kode ke Cabang Utama (*Merge/Pull Request*)

**3. Build (Pembuatan)**
 * [ ] Otomatisasi pemicu pembuatan *build* menggunakan Jenkins (CI)
   * [ ] Mengonfigurasi *Webhook* dari Git ke Jenkins (setiap ada *merge*, Jenkins otomatis jalan)
 * [ ] Membuat Dockerfile
   * [ ] Menggunakan *base image* yang ringan (misal: `node:alpine` atau `node:slim`) untuk menghemat ruang penyimpanan server
   * [ ] Memasukkan perintah instalasi PM2 secara global (`RUN npm install pm2 -g`)
   * [ ] **[NEW]** Menerapkan *Multi-stage Build* (Memisahkan proses instalasi `devDependencies` dan *runtime* agar *image* lebih kecil & aman)
 * [ ] Menyusun *container image* menggunakan Docker
   * [ ] Menjalankan perintah `docker build` lewat Jenkins
   * [ ] **[NEW]** Memanfaatkan *Docker Cache Layer* untuk mempercepat proses *build* di pipeline
 * [ ] Melakukan kompilasi dan pembungkusan kode menjadi *Artifact*
   * [ ] Memastikan tidak ada *error* sintaksis saat proses *build* berjalan

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
   * [ ] Memeriksa apakah ada pustaka npm (*dependencies*) yang memiliki celah keamanan (*vulnerability*) lewat perintah seperti `npm audit` atau integrasi SonarQube

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
