import './Navbar.css'
function Navbar() {
  return (
    <nav className="navbar">
      <div className="logo">React Router</div>
      <ul className="nav-links">
        <li>
          <a href="/">Home</a>
        </li>
        <li>
          <a href="/auth/login">Login</a>
        </li>
      </ul>
    </nav>
  )
}

export default Navbar
