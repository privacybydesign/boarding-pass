import "./App.css";
import Nav from "./components/Nav";
import logoUrl from "./assets/air-logo.svg";
import Router from "./router";

function App() {
  const year = new Date().getFullYear();

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="app-header__brand">
          <img
            src={logoUrl}
            alt="Yivi International Airlines"
            width={56}
            height={56}
          />
          <div className="app-header__copy">
            <h1>Yivi International Airlines</h1>
            <p>Online check-in & boarding pass demo</p>
          </div>
        </div>
        <div className="app-header__meta">
          Check-in opens 24h before departure
        </div>
      </header>
      <Nav />
      <main className="app-main">
        <Router />
      </main>
      <footer
        className="app-footer"
        aria-label="Yivi International Airlines footer"
      >
        <div className="app-footer__inner">
          <span>© {year} Yivi International Airlines</span>
          <span>Privacy-first journeys from check-in to arrival.</span>
        </div>
      </footer>
    </div>
  );
}

export default App;
