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

    // 5. Verifikasi bahwa status berubah menjadi loading atau berhasil (karena backend jalan)
    // Di awal ada kata 'MENUNGGU' atau 'MENJALANKAN...'
    cy.contains('MENJALANKAN...').should('exist');

    // Karena kode default adalah "Halo Dunia" yang benar, kita tunggu sampai berhasil
    cy.contains('DITERIMA (ACCEPTED)!', { timeout: 15000 }).should('be.visible');
  });
});
