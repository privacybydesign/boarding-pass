export default function Nav() {
  return (
    <nav className="main-nav">
      <div className="main-nav__inner">
        <span className="main-nav__label">Online check-in</span>
        <div className="main-nav__links">
          <a href="#/">Home</a>
          <a href="#/buy">Buy</a>
          <a href="#/verify">Verify</a>
          <a href="#/about">About</a>
        </div>
      </div>
    </nav>
  );
}
