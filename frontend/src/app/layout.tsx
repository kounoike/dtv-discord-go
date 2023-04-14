import "./globals.css"

export const metadata = {
  title: "Shichou-Chan",
  description: "視聴ちゃん - Watch with you.",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <head></head>
      <body>{children}</body>
    </html>
  )
}
