import { HashRouter, Route, Routes } from "react-router-dom";
import MainLayout from "./components/layout/MainLayout";
import { routes } from "./routes";

function App() {
    return (
        <HashRouter>
          <Routes>
              <Route path="/" element={<MainLayout/>}>
              {routes}
            </Route>
          </Routes>
        </HashRouter>
      );
}

export default App;
