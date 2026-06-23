# DevOps Best Practices: Alasan di Balik Siklus CI/CD

Dokumen ini disusun sebagai panduan filosofis dan teknis yang menjelaskan **mengapa** setiap tahapan di dalam `devops-cycle.md` itu penting. Mengikuti langkah-langkah ini membedakan seorang pemrogram biasa dari seorang *Software Engineer* kelas *Enterprise*.

---

## 1. Plan & Code (Perencanaan & Pengodean)
**Praktik:** Tidak menulis kode di cabang utama (`main`/`master`).
*   **Mengapa?** Cabang utama adalah cerminan dari aplikasi yang sedang dipakai oleh pengguna secara langsung (*Production*). Jika Anda membuat fitur baru atau mencoba sesuatu dan langsung menyimpannya di sana, aplikasi bisa langsung rusak. Dengan membuat cabang (*branch*) terpisah, Anda memiliki ruang kerja yang aman (*sandbox*).
*   **Best Practice:** Gunakan pola penamaan cabang seperti `feat/login`, `bugfix/typo`, atau `hotfix/db-crash`.

---

## 2. Pre-Commit / Build (Pemeriksaan Awal & Kompilasi)
**Praktik:** Menjalankan Linter (ESLint) dan Security Scanner (GoSec) secara lokal sebelum dikirim ke server.
*   **Mengapa Linter?** Memastikan gaya penulisan kode (*coding style*) konsisten antar developer. Hal ini mencegah perdebatan tidak penting soal spasi atau tanda kutip, dan mendeteksi variabel yang lupa dihapus.
*   **Mengapa GoSec (SAST)?** Mencegah masuknya celah keamanan mendasar (seperti *password* yang ditulis terang-terangan di dalam kode, atau celah *SQL Injection*) **sebelum** kode tersebut terkirim ke internet. Mencegah lebih baik (dan lebih murah) daripada mengobati.
*   **Multi-stage Docker Build:** Kita memisahkan *container* untuk mengompilasi kode dan *container* untuk menjalankan kode. Tujuannya agar ukuran *image* Docker sangat kecil (ringan) dan tidak membawa peralatan *hacker* (seperti *compiler* atau *package manager*) ke server *production*.

---

## 3. Pull Request & Code Review (Tinjauan Kode)
**Praktik:** Kode yang selesai harus melewati mekanisme *Pull Request* (PR) dan diperiksa oleh manusia serta diuji oleh mesin (Jenkins) sebelum digabungkan (*Merge*).
*   **Pemeriksaan Mesin (CI/CD):** Memastikan kode lulus proses kompilasi tanpa campur tangan manusia. "Di laptop saya jalan kok!" adalah alasan yang tidak berlaku lagi jika mesin Jenkins mengatakan gagal.
*   **Pemeriksaan Manusia (Code Review):** Menjaga kualitas logika bisnis dan arsitektur. Mesin bisa tahu kode Anda tidak *error*, tapi hanya manusia (Senior Engineer) yang tahu apakah kode Anda logis dan efisien secara arsitektur.

---

## 4. Container Security Scanning (Trivy)
**Praktik:** Memindai *image* Docker hasil kompilasi dari virus atau *Common Vulnerabilities and Exposures* (CVE).
*   **Mengapa?** Meskipun kode Golang atau Next.js Anda aman, sistem operasi di dalam Docker (seperti Debian atau Alpine) bisa jadi memiliki kelemahan bawaan (misalnya kelemahan pada modul WiFi atau SSL yang belum di-*update* oleh pembuatnya). Trivy memastikan *image* Anda steril dari sisi sistem operasinya.

---

## 5. Automated Testing (Unit Test & Mocking)
**Praktik:** Memastikan fungsi-fungsi kecil berjalan sesuai harapan tanpa bantuan infrastruktur luar (Mocking), serta menetapkan standar *Code Coverage* minimal (misal 80%).
*   **Mengapa Mocking?** *Unit Test* harus sangat cepat (selesai dalam hitungan milidetik). Jika tes memanggil server Docker asli atau MySQL sungguhan, tes akan lambat dan rapuh (*flaky*). Dengan memalsukan (Mock) respon luar, kita mengisolasi pengujian hanya murni pada logika kode buatan kita sendiri.
*   **Mengapa Code Coverage Threshold?** Memaksa disiplin. Jika tim menambahkan 100 baris kode baru tapi lupa menambahkan tes, angka cakupan tes akan turun, dan Jenkins akan menolak kodenya.

---

## 6. Integration Test
**Praktik:** Berbeda dengan *Unit Test*, di sini kita menyalakan MySQL, Redis, dan Docker sungguhan, lalu melihat apakah kode kita bisa berbicara dengan mereka.
*   **Mengapa?** *Mocking* itu bagus untuk kecepatan, tetapi pada akhirnya aplikasi kita benar-benar harus menembak *database* asli. Tes integrasi menjamin bahwa kabel penghubung antarsistem tidak terputus (misalnya: *password* koneksi *database* salah, dsb).

---

## 7. Release & Deploy (Rilis & Peluncuran)
**Praktik:** Memberikan penandaan/label (`git tag` atau versi *Docker Image* seperti `v1.0.1`) sebelum ditarik oleh server *production*.
*   **Mengapa Tagging?** Jika pembaruan aplikasi ternyata menyebabkan server *crash*, DevOps bisa dengan sangat mudah dan instan memundurkan versi (*rollback*) ke label `v1.0.0` sebelumnya, tanpa perlu repot membongkar ulang kode.

---

### Kesimpulan
Secara sekilas, menjalankan seluruh proses `devops-cycle.md` di atas terkesan lambat, merepotkan, dan berlapis-lapis bagi seorang pengerja tunggal (*solo developer*).

Namun, saat aplikasi Anda membesar dan melayani ratusan ribu pengguna harian, siklus kaku inilah yang akan **menyelamatkan Anda** dari *downtime* berjam-jam, peretasan database, atau kepanikan massal akibat aplikasi berhenti bekerja. Anda tidak lagi merilis fitur dengan rasa deg-degan, melainkan dengan ketenangan pikiran karena semua sudah diverifikasi oleh "tentara mesin" Anda.
