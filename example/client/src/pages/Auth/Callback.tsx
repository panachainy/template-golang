import { Navigate } from 'react-router-dom'
import { useAuth } from './AuthContext'
import { useEffect, useState } from 'react'

export function AuthCallbackPage() {
  const { setAccessToken } = useAuth()
  const [error, setError] = useState<string | null>(null)
  const [isProcessing, setIsProcessing] = useState(true)

  useEffect(() => {
    try {
      const params = new URLSearchParams(window.location.search)
      const token = params.get('token')

      if (token) {
        setAccessToken(token)
        localStorage.setItem('jwt', token)
        console.log('Token found in URL:', token)
      } else {
        setError('No authentication token found')
        console.log('No token found in URL')
      }
    } catch (err) {
      setError('Failed to process authentication')
      console.error('Authentication error:', err)
    } finally {
      setIsProcessing(false)
    }
  }, [setAccessToken])

  if (isProcessing) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-100">
        <div className="rounded-lg bg-white p-8 shadow-md">
          <div className="text-center">
            <div className="mb-4 text-gray-600">Processing authentication...</div>
            <div className="h-4 w-24 animate-pulse rounded bg-gray-200" />
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-100">
        <div className="rounded-lg bg-white p-8 shadow-md">
          <div className="text-center">
            <div className="mb-4 text-red-600">{error}</div>
            <Navigate to="/auth/login" replace />
          </div>
        </div>
      </div>
    )
  }

  return <Navigate to="/auth/dashboard" replace />
}
