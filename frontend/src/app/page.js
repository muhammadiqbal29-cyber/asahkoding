"use client";

import { useState, useEffect } from 'react';
import Navbar from './components/Navbar';
import Editor from '@monaco-editor/react';

export default function Home() {
  const [code, setCode] = useState('package main\n\nimport "fmt"\n\nfunc main() {\n  fmt.Println("Halo Dunia")\n}');
  const [output, setOutput] = useState('Konsol siap... Tekan JALANKAN KODE untuk menguji!');
  const [isError, setIsError] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [token, setToken] = useState("");

  // Otomatis mendapatkan token JWT untuk pengujian tanpa membuat UI Login dulu
  useEffect(() => {
    fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: "fauzan@leetcode.com", password: "rahasia123" })
    })
    .then(res => res.json())
    .then(data => {
      if(data.token) setToken(data.token);
    })
    .catch(console.error);
  }, []);

  const handleRunCode = async () => {
    if (!token) {
      setOutput('❌ Error: Token JWT belum didapatkan. Pastikan Golang server menyala.');
      setIsError(true);
      return;
    }

    setIsSubmitting(true);
    setOutput('🚀 Mengirim kode ke Juri Otomatis...\n⏳ Menunggu Docker bekerja di latar belakang...');
    setIsError(false);

    try {
      const response = await fetch('/api/submissions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          problem_id: 1, // ID soal pura-pura
          language: "go",
          code: code
        })
      });

      let data;
      try {
        data = await response.json();
      } catch (e) {
        data = { error: "Response dari server bukan JSON yang valid." };
      }

      if (!response.ok) {
        setOutput(`❌ GAGAL:\n${data.error || 'Terjadi kesalahan sistem.'}`);
        setIsError(true);
      } else {
        if (data.status === "Accepted") {
          setOutput(`✅ DITERIMA (ACCEPTED)!\n\nSeluruh Test Case berhasil dilewati dengan sempurna.\n\nWaktu Eksekusi: ${data.execution_time} ms`);
          setIsError(false);
        } else if (data.status === "Compile Error") {
          setOutput(`❌ GAGAL KOMPILASI (COMPILE ERROR)\n\nCatatan Juri:\n${data.error}`);
          setIsError(true);
        } else {
          setOutput(`❌ ${data.status} (JAWABAN SALAH)\n\nCatatan Juri:\n${data.error_detail}\n\nWaktu Eksekusi: ${data.execution_time} ms`);
          setIsError(true);
        }
      }
    } catch (err) {
      setOutput(`❌ GAGAL KONEKSI:\nBackend Golang tidak terjangkau atau URL salah. Pastikan API menyala di port 8080.\nError: ${err.message}`);
      setIsError(true);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col h-screen">
      <Navbar />

      <div className="flex flex-1 p-6 gap-6 overflow-hidden">
        {/* KIRI: Papan Soal */}
        <div className="flex-[4] flex flex-col overflow-y-auto">
          <div className="gamified-card h-full min-h-full">
            <h1 className="text-3xl font-black mb-2 text-gray-800 tracking-tight">🌳 Penjelajah Kata</h1>
            <div className="flex gap-3 mb-6 font-extrabold text-[15px]">
              <span className="text-brand-green">Mudah</span>
              <span className="text-gray-400">•</span>
              <span className="text-gray-400">Golang</span>
            </div>
            
            <div className="text-base leading-relaxed text-gray-600">
              <p>Selamat datang di AsahKoding! Misimu hari ini adalah menaklukkan level perkenalan ini.</p>
              <br/>
              <p>Tulislah sebuah program yang mencetak kalimat <strong>&quot;Halo Dunia&quot;</strong> ke layar konsol.</p>
              
              <div className="bg-brand-bg rounded-xl border-2 border-gray-200 p-4 mt-4 font-mono text-sm">
                <strong>Contoh 1:</strong><br/><br/>
                <strong>Input:</strong> tidak ada<br/>
                <strong>Ekspektasi Output:</strong> &quot;Halo Dunia&quot;
              </div>
            </div>
          </div>
        </div>

        {/* KANAN: Editor Kode & Konsol Juri */}
        <div className="flex-[6] flex flex-col gap-6">
          <div className="gamified-card flex-1 flex flex-col p-5">
             <div className="flex justify-between items-center mb-4">
               <h3 className="font-extrabold text-gray-400 uppercase text-sm">Editor Kode</h3>
             </div>
             
             <div className="flex-1 rounded-xl overflow-hidden border-2 border-gray-200">
               <Editor
                 height="100%"
                 defaultLanguage="go"
                 value={code}
                 onChange={(val) => setCode(val)}
                 theme="light"
                 options={{
                   minimap: { enabled: false },
                   fontSize: 16,
                   padding: { top: 16 },
                   roundedSelection: true,
                   scrollBeyondLastLine: false,
                 }}
               />
             </div>

             <div className="flex justify-end mt-4">
               <button
                 className={`btn-bouncy px-8 py-4 text-lg ${!token ? 'bg-gray-400 cursor-not-allowed' : 'btn-secondary'}`}
                 onClick={handleRunCode}
                 disabled={isSubmitting || !token}
               >
                 {isSubmitting ? 'MENJALANKAN...' : (!token ? 'MENYIAPKAN...' : '🚀 JALANKAN KODE!')}
               </button>
             </div>
          </div>

          <div className="gamified-card h-[250px] shrink-0 p-5">
             <h3 className="font-extrabold text-gray-400 uppercase text-sm mb-3">Keluaran Konsol</h3>
             <div className={`font-mono bg-[#2b2b2b] p-4 rounded-xl h-[calc(100%-30px)] overflow-y-auto border-4 border-[#1e1e1e] ${isError ? 'text-brand-red' : 'text-brand-green'}`}>
               {output.split('\n').map((line, i) => (
                 <div key={i}>{line}</div>
               ))}
             </div>
          </div>
        </div>
      </div>
    </div>
  );
}
