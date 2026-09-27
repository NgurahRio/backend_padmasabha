import type { Metadata } from "next";
import "./globals.css";
import "./theme.css";

export const metadata: Metadata = { title: "Travora", description: "Travora Manager" };

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>;
}
