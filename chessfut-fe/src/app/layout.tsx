import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL("https://chessfut.com"),
  title: {
    default: "Chessfut — Turn your Chess.com stats into a player card",
    template: "%s | Chessfut",
  },
  description:
    "Turn your Chess.com stats into a FIFA Ultimate Team-style player card, rated out of 99. Compare players and browse the leaderboard.",
  keywords: ["chess", "chess.com", "player card", "FIDE", "chess rating", "chess stats"],
  openGraph: {
    type: "website",
    siteName: "Chessfut",
    title: "Chessfut — Turn your Chess.com stats into a player card",
    description:
      "Turn your Chess.com stats into a FIFA Ultimate Team-style player card, rated out of 99.",
    url: "https://chessfut.com",
  },
  twitter: {
    card: "summary_large_image",
    title: "Chessfut — Turn your Chess.com stats into a player card",
    description:
      "Turn your Chess.com stats into a FIFA Ultimate Team-style player card, rated out of 99.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}