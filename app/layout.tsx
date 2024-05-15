import type { Metadata } from "next";
import { Fira_Mono } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";

const FiraMonoFont = Fira_Mono({
  subsets: ["latin"],
  weight: ["400", "500", "700"],
});

export const metadata: Metadata = {
  title: "YT Downloader",
  description: "Download YouTube videos",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={cn(
          FiraMonoFont.className,
          "flex items-center justify-center w-full h-screen"
        )}
      >
        <div className="fixed -z-50 h-screen w-full animate-pulse bg-gradient-to-br from-blue-400 to-pink-400" />
        {children}
      </body>
    </html>
  );
}
