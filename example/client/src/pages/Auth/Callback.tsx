import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

// FIXME: check navigate path
export function AuthCallbackPage() {
  const navigate = useNavigate()

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const token = params.get('token')

    if (token) {
      localStorage.setItem('jwt', token)
      navigate('auth/dashboard')
    } else {
      navigate('auth/login')
    }
  }, [navigate])

  return <div>Processing authentication...</div>
}
