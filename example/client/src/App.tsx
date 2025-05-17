import './App.css'

import { LogsProvider } from '@/providers/Logs'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import ErrorPage from './core/Error'
import { MainLayout } from './layouts/Main'
import { AuthProvider } from './pages/Auth/AuthContext'
import { AuthCallbackPage } from './pages/Auth/Callback'
import { LoginPage } from './pages/Auth/Login'
import Home from './pages/Home/Home'

const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <MainLayout>
        <Home />
      </MainLayout>
    ),
    errorElement: <ErrorPage />,
  },
  {
    path: 'auth/login',
    element: (
      <MainLayout>
        <LoginPage />
      </MainLayout>
    ),
    errorElement: <ErrorPage />,
  },
  {
    path: 'auth/callback',
    element: (
      <MainLayout>
        <AuthCallbackPage />
      </MainLayout>
    ),
    errorElement: <ErrorPage />,
  },
])

function App() {
  return (
    <AuthProvider>
      <LogsProvider>
        <RouterProvider router={router} />
      </LogsProvider>
    </AuthProvider>
  )
}

export default App
