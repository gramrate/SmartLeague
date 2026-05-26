import { Link, useNavigate, useRouterState } from "@tanstack/react-router";
import { useAuthStore } from "@/lib/auth-store";
import { displayUserName } from "@/lib/roles";
import { Button } from "@/components/ui/button";
import { Trophy, LogOut, User as UserIcon, Menu, X } from "lucide-react";
import { useState } from "react";
import { cn } from "@/lib/utils";

export function Header() {
  const me = useAuthStore((s) => s.me);
  const logout = useAuthStore((s) => s.logout);
  const navigate = useNavigate();
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const [open, setOpen] = useState(false);
  const nav = [
    { to: "/series", label: "Серии" },
    { to: "/clubs", label: "Клубы" },
    { to: "/players", label: "Игроки" },
    ...(me?.club_id ? [{ to: "/clubs/$id" as const, label: "Мой клуб", params: { id: me.club_id } }] : []),
  ];
  const isNavActive = (item: (typeof nav)[number]) => {
    if ("params" in item && item.params?.id) {
      return pathname.startsWith(`/clubs/${item.params.id}`);
    }
    if (item.to === "/clubs" && me?.club_id && pathname.startsWith(`/clubs/${me.club_id}`)) {
      return false;
    }
    return pathname.startsWith(item.to);
  };

  const handleLogout = async () => {
    await logout();
    navigate({ to: "/" });
  };

  return (
    <header className="sticky top-0 z-40 border-b border-border/60 bg-background/80 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6">
        <Link to="/" className="flex items-center gap-2 font-display text-lg font-bold tracking-tight">
          <span className="grid h-8 w-8 place-items-center rounded-md bg-gradient-to-br from-primary to-primary-glow shadow-[var(--shadow-glow)]">
            <Trophy className="h-4 w-4 text-primary-foreground" />
          </span>
          <span>SmartLeague</span>
        </Link>

        <nav className="hidden items-center gap-1 md:flex">
          {nav.map((n) => (
            <Link
              key={`${n.to}-${n.label}`}
              to={n.to}
              params={("params" in n ? n.params : undefined) as never}
              className={cn(
                "rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground",
                isNavActive(n) && "bg-secondary text-foreground",
              )}
            >
              {n.label}
            </Link>
          ))}
        </nav>

        <div className="hidden items-center gap-2 md:flex">
          {me ? (
            <>
              <Link to="/account" className="flex items-center gap-2 rounded-md px-3 py-2 text-sm text-muted-foreground hover:text-foreground">
                <UserIcon className="h-4 w-4" />
                {displayUserName(me)}
              </Link>
              <Button variant="ghost" size="sm" onClick={handleLogout}>
                <LogOut className="h-4 w-4" />
              </Button>
            </>
          ) : (
            <>
              <Button variant="ghost" size="sm" asChild>
                <Link to="/login">Вход</Link>
              </Button>
              <Button size="sm" asChild>
                <Link to="/register">Регистрация</Link>
              </Button>
            </>
          )}
        </div>

        <button
          className="rounded-md p-2 text-muted-foreground hover:text-foreground md:hidden"
          onClick={() => setOpen((v) => !v)}
          aria-label="Меню"
        >
          {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>
      </div>

      {open && (
        <div className="border-t border-border/60 md:hidden">
          <div className="mx-auto flex max-w-7xl flex-col px-4 py-3">
            {nav.map((n) => (
              <Link key={`${n.to}-${n.label}`} to={n.to} params={("params" in n ? n.params : undefined) as never} onClick={() => setOpen(false)}
                className="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-secondary hover:text-foreground">
                {n.label}
              </Link>
            ))}
            <div className="my-2 border-t border-border/60" />
            {me ? (
              <>
                <Link to="/account" onClick={() => setOpen(false)} className="rounded-md px-3 py-2 text-sm text-muted-foreground hover:text-foreground">
                  Аккаунт ({displayUserName(me)})
                </Link>
                <button onClick={handleLogout} className="rounded-md px-3 py-2 text-left text-sm text-muted-foreground hover:text-foreground">
                  Выйти
                </button>
              </>
            ) : (
              <>
                <Link to="/login" onClick={() => setOpen(false)} className="rounded-md px-3 py-2 text-sm hover:text-foreground">Вход</Link>
                <Link to="/register" onClick={() => setOpen(false)} className="rounded-md px-3 py-2 text-sm font-medium text-primary">Регистрация</Link>
              </>
            )}
          </div>
        </div>
      )}
    </header>
  );
}
