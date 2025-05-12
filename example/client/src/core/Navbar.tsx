import './Navbar.css' // นำเข้า Navbar.css สำหรับการตกแต่ง

function Navbar() {
  return (
    <nav className="navbar">
      <div className="logo">React Router</div>
      <ul className="nav-links">
        <li>
          <a href="/">Home</a>
        </li>
        <li>
          <a href="/pokemon">Pokemon</a>
        </li>
        <li>
          <a href="/zod">Zod</a>
        </li>
        <li>
          <a href="/styled">Styled</a>
        </li>
      </ul>
    </nav>
  )
}

export default Navbar
