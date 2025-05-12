import './App.css'

import { LogsProvider } from '@/providers/Logs'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import ErrorPage from './core/Error'
import { MainLayout } from './layouts/Main'
import { CallbackPage } from './pages/Auth/Callback'
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
        <CallbackPage />
      </MainLayout>
    ),
    errorElement: <ErrorPage />,
    loader: ({ request }) => {
      const url = new URL(request.url)
      const token = url.searchParams.get('token')
      if (token) {
        console.log('Token:', token)
        localStorage.setItem('token', token)
      }
      return null
    },
  },
])

function App() {
  return (
    <LogsProvider>
      <RouterProvider router={router} />
    </LogsProvider>
  )
}

export default App
