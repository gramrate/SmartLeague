import { createBrowserRouter, Navigate } from "react-router-dom";
import { RootLayout } from "./rootLayout";
import { ErrorPage } from "../screens/system/errorPage";
import { LoginPage } from "../screens/auth/login";
import { RegisterPage } from "../screens/auth/register";
import { ClubsPage } from "../screens/clubs/list";
import { ClubDetailPage } from "../screens/clubs/detail";
import { ClubGamesPage } from "../screens/clubs/games";
import { ClubManagePage } from "../screens/clubs/manage";
import { ClubCreatePage } from "../screens/clubs/create";
import { SeriesDetailPage } from "../screens/series/detail";
import { SeriesCreatePage } from "../screens/series/create";
import { ClubSeriesPage } from "../screens/series/byClub";
import { AllSeriesPage } from "../screens/series/all";
import { GameDetailPage } from "../screens/games/detail";
import { GameManagePage } from "../screens/games/manage";
import { AccountPage } from "../screens/user/account";
import { UserViewPage } from "../screens/user/view";
import { UserGamesPage } from "../screens/user/games";
import { UserSeriesPage } from "../screens/user/series";
import { PlayersPage } from "../screens/players/list";
import { ProtectedRoute } from "./protectedRoute";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <Navigate to="/clubs" replace /> },
      { path: "login", element: <LoginPage /> },
      { path: "register", element: <RegisterPage /> },
      { path: "players", element: <PlayersPage /> },
      { path: "user/:id", element: <UserViewPage /> },
      { path: "user/:id/games", element: <UserGamesPage /> },
      { path: "user/:id/series", element: <UserSeriesPage /> },
      { path: "clubs", element: <ClubsPage /> },
      { path: "clubs/:id", element: <ClubDetailPage /> },
      { path: "clubs/:id/manage", element: <ClubManagePage /> },
      { path: "clubs/:id/games", element: <ClubGamesPage /> },
      { path: "clubs/:id/series", element: <ClubSeriesPage /> },
      { path: "series", element: <AllSeriesPage /> },
      { path: "series/:id", element: <SeriesDetailPage /> },
      { path: "game/:id", element: <GameDetailPage /> },
      {
        element: <ProtectedRoute />,
        children: [
          { path: "account", element: <AccountPage /> },
          { path: "clubs/create", element: <ClubCreatePage /> },
          { path: "series/create", element: <SeriesCreatePage /> },
          { path: "game/:id/manage", element: <GameManagePage /> }
        ]
      }
    ]
  }
]);
