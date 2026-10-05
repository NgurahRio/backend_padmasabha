import type { Metadata } from "next";
import "./globals.css";
import "./theme.css";

export const metadata: Metadata = {
  title: "Padmanabha",
  description: "Padmanabha Manager",
  icons: { icon: "/padmanabha-logo.png" },
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>;
}
