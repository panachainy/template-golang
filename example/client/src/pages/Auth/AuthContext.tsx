import { useContext } from 'react'
import { createContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import type { AuthContextType } from './interfaces/AuthContext'
import type { UserInfo } from './interfaces/UserInfo'

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null)

  // Check authentication status on mount and after token refresh
  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        const response = await fetch(
          'http://localhost:8080/api/v1/auth/status',
          {
            credentials: 'include', // This is important for sending cookies
          },
        )
        if (response.ok) {
          const data = await response.json()
          setUserInfo({
            accessToken: data.accessToken || null,
            refreshToken: data.refreshToken || null,
          })
        } else {
          setUserInfo(null)
        }
      } catch (error) {
        console.error('Failed to check auth status:', error)
        setUserInfo(null)
      }
    }

    checkAuthStatus()
  }, [])

  // Function to refresh access token
  const refreshAccessToken = async () => {
    console.log('TODO: implement later Refreshing access token...')
  }

  // Function to handle login
  const loginWithLine = async () => {
    window.location.href = 'http://localhost:8080/api/v1/auth/line/login'
  }

  const setAccessToken = (accessToken: string | null) => {
    setUserInfo((prevUserInfo) => {
      const newUserInfo = {
        accessToken,
        refreshToken: prevUserInfo?.refreshToken || '',
      }
      return newUserInfo
    })
  }

  return (
    <AuthContext.Provider
      value={{ userInfo, setAccessToken, loginWithLine, refreshAccessToken }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const PrivateRoute = ({ children }: { children: ReactNode }) => {
  const authCtx = useContext(AuthContext)

  if (!authCtx) {
    throw new Error('PrivateRoute must be used within an AuthProvider')
  }

  return authCtx.userInfo ? children : <Navigate to="auth/login" />
}

const UseAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

export { AuthContext, UseAuth }
