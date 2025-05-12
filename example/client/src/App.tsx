import './App.css'

import { LogsProvider } from '@/providers/Logs'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import ErrorPage from './core/Error'
import { MainLayout } from './layouts/Main'
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
    path: '/login',
    element: (
      <MainLayout>
        <LoginPage />
      </MainLayout>
    ),
    errorElement: <ErrorPage />,
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
