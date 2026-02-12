import './globals.css'

export const metadata = {
  title: 'ClipWave AI - Transform Videos Into Viral Clips',
  description: 'AI-powered video clipping platform that automatically detects and clips the best moments from your YouTube videos',
}

export default function RootLayout({ children }) {
  return (
    <html lang="en" className="dark">
      <body className="antialiased">
        {children}
      </body>
    </html>
  )
}
