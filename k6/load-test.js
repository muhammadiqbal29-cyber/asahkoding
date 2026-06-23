import http from 'k6/http';
import { check, sleep } from 'k6';

// Konfigurasi Load Test
export let options = {
    // Stage 1: Naik ke 50 pengguna dalam 2 detik
    // Stage 2: Tahan 50 pengguna selama 8 detik
    // Stage 3: Turun ke 0 pengguna dalam 2 detik
    stages: [
        { duration: '2s', target: 50 },
        { duration: '8s', target: 50 },
        { duration: '2s', target: 0 },
    ],
    // Ambang batas kelulusan (Thresholds)
    thresholds: {
        // Kegagalan request harus kurang dari 5%
        http_req_failed: ['rate<0.05'], 
        // 95% request harus selesai di bawah 500ms
        http_req_duration: ['p(95)<500'], 
    },
};

// URL Backend (Karena dijalankan di dalam Docker Network)
const BASE_URL = 'http://backend-test:8080';

export default function () {
    // Menguji Endpoint Utama
    let res = http.get(`${BASE_URL}/health`);
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    // Menguji Endpoint Redis Integration
    let resRedis = http.get(`${BASE_URL}/health/redis`);
    check(resRedis, {
        'redis is connected': (r) => r.status === 200,
    });

    // Beri jeda sedikit agar tidak murni spam tanpa henti (mirip perilaku user asli)
    sleep(0.1);
}
