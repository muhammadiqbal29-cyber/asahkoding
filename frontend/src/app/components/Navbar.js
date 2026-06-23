export default function Navbar() {
  return (
    <nav className="flex justify-between items-center px-8 py-4 bg-white border-b-2 border-gray-200 sticky top-0 z-50">
      <div className="text-2xl font-black text-brand-green tracking-tight">
        AsahKoding 🐒
      </div>
      <div className="flex gap-6">
        <a href="#" className="text-brand-blue font-bold text-[15px] uppercase transition-colors hover:text-brand-blue">Latihan</a>
        <a href="#" className="text-gray-400 font-bold text-[15px] uppercase transition-colors hover:text-brand-blue">Materi</a>
        <a href="#" className="text-gray-400 font-bold text-[15px] uppercase transition-colors hover:text-brand-blue">Papan Skor</a>
      </div>
      <div className="flex gap-4">
        <div className="flex items-center gap-2 px-4 py-2 rounded-full font-extrabold text-sm border-2 border-gray-200 bg-gray-50 text-gray-700">
          💎 1,450
        </div>
        <div className="flex items-center gap-2 px-4 py-2 rounded-full font-extrabold text-sm border-2 border-[#ffd8a8] bg-[#fff4e5] text-[#e67700]">
          🔥 32
        </div>
      </div>
    </nav>
  );
}
