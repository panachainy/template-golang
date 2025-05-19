import { useAuth } from '@/pages/Auth/AuthContext'
import './Navbar.css'
function Navbar() {
  const { userInfo } = useAuth()

  return (
    <nav className="navbar">
      <div className="logo">React Router</div>
      <div className="user-status">
        {userInfo ? `Logged in as ${userInfo.accessToken}` : 'Not logged in'}
      </div>
      <ul className="nav-links">
        <li>
          <a href="/">Home</a>
        </li>
        <li>
          <a href="/auth/login">Login</a>
        </li>
        <li>
          <a href="/auth/dashboard">Dashboard</a>
        </li>
      </ul>
    </nav>
  )
}

export default Navbar
