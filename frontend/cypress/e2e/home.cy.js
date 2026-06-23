describe('AsahKoding E2E Tests', () => {
  it('should load the homepage and display the editor', () => {
    // Abaikan error dari third-party library (Monaco Editor CDN, dsb.)
    cy.on('uncaught:exception', () => false);

    // 1. Kunjungi halaman utama
    cy.visit('/', { timeout: 30000 });

    // 2. Pastikan judul aplikasi muncul (membuktikan halaman ter-render)
    cy.contains('Penjelajah Kata', { timeout: 30000 }).should('be.visible');

    // 3. Tunggu tombol JALANKAN KODE muncul (membuktikan useEffect sudah jalan & JWT didapat)
    //    Ini lebih robust daripada cy.wait('@autoLogin') karena langsung mengecek hasil akhir.
    cy.contains('JALANKAN KODE!', { timeout: 30000 }).should('be.visible');

    // 4. Klik tombol jalankan kode
    cy.contains('JALANKAN KODE!').click();

    // 5. Tunggu hasil eksekusi muncul (Accepted)
    //    Backend perlu waktu untuk compile & run kode Go di Docker container.
    cy.contains('DITERIMA (ACCEPTED)!', { timeout: 60000 }).should('be.visible');
  });
});
