# 🚀 Lean DevOps: Arsitektur CI/CD Startup

Dokumen ini mendeskripsikan siklus pengembangan dan CI/CD yang telah dioptimalkan secara ekstrem untuk *Solo Developer* dan *Startup* tahap awal. Filosofi utamanya adalah **Fast Feedback Loop**: jangan habiskan waktu menunggu *pipeline* berjalan, habiskan waktu untuk *coding*.

Berbeda dengan `devops-cycle.md` (yang merupakan rancangan kelas Enterprise), *pipeline* ini dirancang pragmatis: memisahkan *feedback* harian dari pengujian skala besar.

---

## 🛤️ Dua Jalur Pipeline

Pipeline Jenkins kita (diatur dalam `Jenkinsfile`) beroperasi dalam dua mode yang berbeda tergantung pada seberapa penting *commit* Anda:

### 1. 🟢 Standard Pipeline (Jalur Kilat)
**Pemicu:** Push rutin ke *branch* `main`.
**Estimasi Waktu:** 1 - 2 Menit.
**Tujuan:** Memastikan aplikasi bisa dikompilasi, lolos standar kode (*lint*), *unit test* sukses, dan berhasil dibungkus menjadi *Docker Image* siap rilis.
**Tahapan:**
1. **Checkout SCM**: Menarik kode terbaru dari GitHub.
2. **Pre-Commit Checks**: ESLint & GoSec (berjalan cepat secara lokal sebelum *commit*, atau divalidasi cepat di Jenkins).
3. **Build Docker Images**: Mengemas aplikasi menggunakan *multi-stage build* dan *cache*.
4. **Unit Test**: Menjalankan pengujian per fungsi (`go test`) untuk memastikan logika inti aman.
5. **Release & Semantic Versioning**: Melakukan penomoran rilis (v1.x.x) otomatis menggunakan *Conventional Commits*, melakukan *tagging* di GitHub, dan mendorong *Docker Image* terbaru ke Docker Hub.

*Catatan: Segala jenis pengujian yang melibatkan "Docker-in-Docker" atau proses yang makan waktu ditiadakan dari jalur ini untuk menghemat memori laptop.*

---

### 2. 🔴 Heavy Pipeline (Jalur Rilis Besar / Nightly)
**Pemicu:** Eksekusi manual dari *dashboard* Jenkins dengan mencentang parameter `RUN_HEAVY_TESTS`.
**Estimasi Waktu:** 10 - 15 Menit.
**Tujuan:** Memastikan aplikasi tidak hanya berfungsi secara logika, tetapi juga lolos audit keamanan secara penuh dan lulus pengujian secara sistem-ke-sistem (*End-to-End*).
**Tahapan (Menambahkan tahap berikut pada Standard Pipeline):**
1. **Security Scan (Trivy)**: Memindai kerentanan OS pada *Docker Image*.
2. **Dependency Audit**: Menjalankan `npm audit` dan `govulncheck` untuk mencari CVE pada *library* pihak ketiga.
3. **Docker Compose E2E Setup**: Jenkins akan mengangkat miniatur *production* (MySQL, Redis, Backend, Frontend) di dalam jaringannya secara paralel.
4. **Integration Test (cURL)**: Memastikan Backend dapat berkomunikasi mulus dengan MySQL dan Redis.
5. **Load & Stress Test (K6)**: Memberikan beban ribuan permintaan simulasi per detik.
6. **DAST Scan (OWASP ZAP)**: Mensimulasikan serangan *hacker* secara dinamis ke *endpoint* aplikasi yang sedang menyala.
7. **End-to-End Test (Cypress)**: Menjalankan *browser headless* (Electron/Chrome) yang melakukan klik dan simulasi pengetikan persis seperti pengguna manusia (Membuka halaman, login, jalankan kode, dapat *Accepted*).

---

## 🛠️ Manajemen Perubahan Harian (Workflow)

1. Tulis kode fitur Anda (Frontend/Backend).
2. Tulis *Unit Test* kecil untuk fungsi tersebut.
3. Lakukan komit standar, misalnya: `git commit -m "feat: tambah halaman soal"`.
4. *Push* ke repositori: `git push origin main`.
5. Anda hanya perlu menunggu ~1 menit. Jika *pipeline* berwarna hijau, fitur Anda sudah terbungkus rapi di *Docker Hub*.
6. Jika Anda sudah mengumpulkan beberapa fitur dan bersiap melakukan peluncuran (*Launch*), buka Jenkins, klik *Build with Parameters*, centang `RUN_HEAVY_TESTS`, dan klik *Build*. Sambil menunggu 15 menit, Anda bisa pergi menyeduh kopi! ☕
