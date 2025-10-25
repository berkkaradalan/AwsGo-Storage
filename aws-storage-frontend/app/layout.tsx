import "./globals.css"

export const metadata = {
  title: "AWS Storage Dashboard",
  description: "Demo Dashboard",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>{children}</body>
    </html>
  )
}