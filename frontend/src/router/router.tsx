import { createBrowserRouter, Navigate } from "react-router-dom";
import { RootLayout } from "./rootLayout";
import { ErrorPage } from "../screens/system/errorPage";
import { LoginPage } from "../screens/auth/login";
import { RegisterPage } from "../screens/auth/register";
import { ClubsPage } from "../screens/clubs/list";
import { ClubDetailPage } from "../screens/clubs/detail";
import { ClubCreatePage } from "../screens/clubs/create";
import { SeriesDetailPage } from "../screens/series/detail";
import { SeriesCreatePage } from "../screens/series/create";
import { ClubSeriesPage } from "../screens/series/byClub";
import { GameDetailPage } from "../screens/games/detail";
import { GameManagePage } from "../screens/games/manage";
import { AccountPage } from "../screens/user/account";
import { UserViewPage } from "../screens/user/view";
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
      {
        element: <ProtectedRoute />,
        children: [
          { path: "account", element: <AccountPage /> },
          { path: "user/:id", element: <UserViewPage /> },
          { path: "clubs", element: <ClubsPage /> },
          { path: "clubs/create", element: <ClubCreatePage /> },
          { path: "clubs/:id", element: <ClubDetailPage /> },
          { path: "clubs/:id/series", element: <ClubSeriesPage /> },
          { path: "series/:id", element: <SeriesDetailPage /> },
          { path: "series/create", element: <SeriesCreatePage /> },
          { path: "game/:id", element: <GameDetailPage /> },
          { path: "game/:id/manage", element: <GameManagePage /> }
        ]
      }
    ]
  }
]);
