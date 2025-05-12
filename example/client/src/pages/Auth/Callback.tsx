import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

export function AuthCallbackPage() {
  const navigate = useNavigate()

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const token = params.get('token')

    if (token) {
      localStorage.setItem('jwt', token)
      navigate('/dashboard')
    } else {
      navigate('/login')
    }
  }, [navigate])

  return <div>Processing authentication...</div>
}
