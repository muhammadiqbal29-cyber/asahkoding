import { NextResponse } from 'next/server';

export function middleware(request) {
  // Hanya proses jika path dimulai dengan /api/
  if (request.nextUrl.pathname.startsWith('/api/')) {
    // Ambil BACKEND_URL dari environment variable saat RUNTIME (saat kontainer berjalan)
    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';
    
    // Ganti localhost:3000/api/auth/login menjadi http://backend-test:8080/api/auth/login
    const destination = request.nextUrl.href.replace(request.nextUrl.origin, backendUrl);
    
    // Rewrite akan membuat Next.js bertindak sebagai Reverse Proxy (BFF)
    return NextResponse.rewrite(new URL(destination));
  }
}

export const config = {
  // Hanya jalankan middleware ini untuk endpoint API
  matcher: '/api/:path*',
};
