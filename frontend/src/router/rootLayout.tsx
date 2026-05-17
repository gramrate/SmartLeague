import { Link, NavLink, Outlet } from "react-router-dom";
import { useAuthStore } from "../store/authStore";
import { logoutUser } from "../api/auth";

function cn({ isActive }: { isActive: boolean }) {
  return isActive ? "text-blue-600 font-medium" : "text-gray-700";
}

export function RootLayout() {
  const { isAuthed, clear } = useAuthStore();

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      <header className="border-b bg-white">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
          <Link to="/" className="font-semibold">
            SmartLeague
          </Link>
          <nav className="flex gap-4 text-sm">
            <NavLink to="/clubs" className={cn}>
              Clubs
            </NavLink>
            <NavLink to="/series" className={cn}>
              Series
            </NavLink>
            <NavLink to="/players" className={cn}>
              Players
            </NavLink>
            <NavLink to="/account" className={cn}>
              Account
            </NavLink>
          </nav>
          <div className="flex items-center gap-3">
            {isAuthed ? (
              <button
                className="rounded bg-gray-900 px-3 py-1.5 text-sm text-white"
                onClick={async () => {
                  try {
                    await logoutUser();
                  } finally {
                    clear();
                  }
                }}
              >
                Logout
              </button>
            ) : (
              <div className="flex gap-2">
                <Link className="text-sm text-blue-600" to="/login">
                  Login
                </Link>
                <Link className="text-sm text-blue-600" to="/register">
                  Register
                </Link>
              </div>
            )}
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
