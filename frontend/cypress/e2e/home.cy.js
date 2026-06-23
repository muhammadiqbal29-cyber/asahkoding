describe('AsahKoding E2E Tests', () => {
  it('should load the homepage and display the editor', () => {
    // 1. Kunjungi halaman utama
    cy.visit('/');

    // 2. Pastikan judul aplikasi muncul
    cy.contains('Penjelajah Kata').should('be.visible');

    // 3. Pastikan tombol jalankan kode ada
    cy.contains('JALANKAN KODE!').should('be.visible');

    // 4. Klik tombol jalankan kode
    cy.contains('JALANKAN KODE!').click();

    // 5. Kita langsung tunggu hasil akhirnya saja, karena kecepatan eksekusi API di backend
    // kadang kurang dari 50ms sehingga tulisan 'MENJALANKAN...' lewat terlalu cepat.
    // Karena kode default adalah "Halo Dunia" yang benar, kita tunggu sampai berhasil
    cy.contains('DITERIMA (ACCEPTED)!', { timeout: 15000 }).should('be.visible');
  });
});
