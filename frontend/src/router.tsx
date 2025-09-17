import { useEffect, useState } from "react";
import Home from "./pages/Home";
import About from "./pages/About";
import BuyTicket from "./pages/BuyTicket";
import Verify from "./pages/Verify";

type Route = "#/" | "#/about" | "#/buy" | "#/verify";

function getRoute(): Route {
  const h = window.location.hash || "#/";
  switch (h) {
    case "#/about":
    case "#/buy":
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
    case "#/buy":
      return <BuyTicket />;
    case "#/verify":
      return <Verify />;
    case "#/":
    default:
      return <Home />;
  }
}
