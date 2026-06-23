import { Nunito } from "next/font/google";
import "./globals.css";

// Menggunakan font membulat (rounded) yang sangat cocok untuk UI Gamified
const nunito = Nunito({
  variable: "--font-nunito",
  subsets: ["latin"],
  weight: ["400", "600", "700", "800", "900"],
});

export const metadata = {
  title: "AsahKoding",
  description: "Asah logika kodemu layaknya bermain game!",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body className={`${nunito.variable}`}>
        {children}
      </body>
    </html>
  );
}
