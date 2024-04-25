import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Navigation } from "@/components/nav";
import { Providers } from "./providers";
import { Toaster } from "@/components/ui/sonner";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Directory Uploader",
  description: "Manage a contact list uploaded from CSV",
  icons: {
    icon: [
      { url: '/favicon.ico' }
    ]
  }
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Providers>
          <div className="flex border-b-current border-2 items-center px-8">
            <p className="font-mono font-bold">
              📧 Email Uploader 📒
            </p>
            <Navigation />
          </div>
          {children}
          <Toaster />
        </Providers>
      </body>
    </html>
  );
}
