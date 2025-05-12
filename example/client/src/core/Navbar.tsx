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
          <a href="/login">Login</a>
        </li>
      </ul>
    </nav>
  )
}

export default Navbar
