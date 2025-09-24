import { useEffect, useState } from "react";
import Home from "./pages/Home";
import About from "./pages/About";
import Verify from "./pages/Verify";

type Route = "#/" | "#/about" | "#/verify";

function getRoute(): Route {
  const h = window.location.hash || "#/";
  switch (h) {
    case "#/about":
    case "#/verify":
      return h as Route;
    default:
      return "#/";
  }
}

export default function Router() {
  const [route, setRoute] = useState<Route>(getRoute());

  useEffect(() => {
    const onHashChange = () => setRoute(getRoute());
    window.addEventListener("hashchange", onHashChange);
    return () => window.removeEventListener("hashchange", onHashChange);
  }, []);

  switch (route) {
    case "#/about":
      return <About />;
    case "#/verify":
      return <Verify />;
    case "#/":
    default:
      return <Home />;
  }
}
