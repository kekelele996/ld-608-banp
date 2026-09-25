import { useState } from "react";
import { createRoot } from "react-dom/client";
import { routes } from "./router/routes";
import { ReleasePage } from "./pages/ReleasePage";
import { PlaceholderPage } from "./pages/PlaceholderPage";
import "./styles.css";

function App() {
  const [active, setActive] = useState<string>(routes[0]?.route ?? "/release");
  const current = routes.find((route) => route.route === active) ?? routes[0];

  return (
    <div className="shell">
      <aside>
        <div className="brand">航空地勤周转保障平台</div>
        <nav>
          {routes.map((route) => (
            <button
              key={route.route}
              className={active === route.route ? "active" : ""}
              onClick={() => setActive(route.route)}
            >
              {route.name}
            </button>
          ))}
        </nav>
      </aside>
      {active === "/release" ? (
        <ReleasePage />
      ) : (
        <PlaceholderPage name={current?.name ?? "工作台"} />
      )}
    </div>
  );
}

createRoot(document.getElementById("root")!).render(<App />);
