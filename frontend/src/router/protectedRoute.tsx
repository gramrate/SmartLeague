import { useEffect } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuthStore } from "../store/authStore";
import { Role } from "../types/enums";

export function ProtectedRoute({ minRole }: { minRole?: Role }) {
  const { isReady, isAuthed, init, role } = useAuthStore();
  const location = useLocation();

  useEffect(() => {
    if (!isReady) void init();
  }, [isReady, init]);

  if (!isReady) {
    return <div className="p-6">Loading...</div>;
  }

  if (!isAuthed) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  if (minRole != null && (role ?? Role.User) < minRole) {
    return <div className="p-6">403: forbidden</div>;
  }

  return <Outlet />;
}

