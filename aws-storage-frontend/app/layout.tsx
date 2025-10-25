// import "./globals.css"

// export const metadata = {
//   title: "AWS Storage Dashboard",
//   description: "Demo Dashboard",
// }

// export default function RootLayout({ children }: { children: React.ReactNode }) {
//   return (
//     <html lang="en" suppressHydrationWarning>
//       <body>{children}</body>
//     </html>
//   )
// }


import "./globals.css"
import { Providers } from "./providers"

export const metadata = {
  title: "AWS Storage Dashboard",
  description: "Demo Dashboard",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  )
}