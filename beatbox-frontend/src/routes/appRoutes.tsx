import SearchPage from "../pages/dashboard/SearchPage";
import HomePage from "../pages/home/HomePage";
import { RouteType } from "./config";
import SearchIcon from '@mui/icons-material/Search';
import HomeIcon from '@mui/icons-material/Home';

const appRoutes: RouteType[] = [
  {
    index: true,
    element: <HomePage />,
    state: "home"
  },
  {
    path: "/",
    element: <HomePage />,
    state: "Home",
    sidebarProps: {
      displayText: "Home",
      icon: <HomeIcon />
    }
  },
  { //FIXME: Replace with search
      path: "/search",
      element: <SearchPage />,
      state: "search",
      sidebarProps: {
          displayText: "Search",
          icon: <SearchIcon />
      },
  },
];

export default appRoutes;