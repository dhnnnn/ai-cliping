import { NextResponse } from 'next/server'

export async function GET(request) {
  const { pathname } = new URL(request.url)
  
  // Health check endpoint
  if (pathname === '/api/' || pathname === '/api') {
    return NextResponse.json({ 
      message: 'ClipWave AI Backend Running',
      status: 'healthy',
      timestamp: new Date().toISOString()
    })
  }

  return NextResponse.json({ error: 'Not Found' }, { status: 404 })
}

export async function POST(request) {
  const { pathname } = new URL(request.url)
  
  // Placeholder for future API endpoints
  return NextResponse.json({ error: 'Endpoint not implemented' }, { status: 404 })
}
