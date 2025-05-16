import { useContext } from 'react'
import { Navigate } from 'react-router-dom'

import { createContext, useState } from 'react'
import type { ReactNode } from 'react'
import type { UserInfo } from './interfaces/UserInfo'

interface AuthContextType {
  userInfo: unknown | null
  setAccessToken: (accessToken: string) => void

  loginWithLine: () => void
  refreshAccessToken: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [userInfo, setUserInfo] = useState<UserInfo | null>()

  // Function to refresh access token
  const refreshAccessToken = async () => {
    // TODO: implement later

    console.log('TODO: implement later Refreshing access token...')
    // try {
    //   const response = await axios.post(
    //     '/api/refresh-token',
    //     {},
    //     { withCredentials: true },
    //   )
    //   setAccessToken(response.data.accessToken)
    // } catch (error) {
    //   console.error('Failed to refresh access token:', error)
    //   setAccessToken(null)
    // }
  }

  // Function to handle login
  const loginWithLine = async () => {
    console.log('TODO: implement later Logging in with LINE...')

    window.location.href = 'http://localhost:8080/api/v1/auth/line/login'
    // try {
    //   const response = await axios.get('/api/sso-login', { withCredentials: true });
    //   setAccessToken(response.data.accessToken);
    // } catch (error) {
    //   console.error('SSO login failed:', error);
    // }
  }

  const setAccessToken = (accessToken: string | null) => {
    setUserInfo((prevUserInfo) => ({
      accessToken,
      refreshToken: prevUserInfo?.refreshToken || '',
    }))
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

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
