import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { usersApi, clubsApi, ApiError } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { RoleBadge } from "@/components/site/RoleBadge";
import { ClubState } from "@/types/api";
import { displayUserName } from "@/lib/roles";
import { fmtDate, fmtDateRange, fmtRub } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/lib/auth-store";
import { toast } from "sonner";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";

export const Route = createFileRoute("/user/$id")({ component: UserPage });

function UserPage() {
  const { id } = Route.useParams();
  const location = useLocation();
  const me = useAuthStore((s) => s.me);
  const qc = useQueryClient();
  const user = useQuery({ queryKey: ["user", id], queryFn: () => usersApi.get(id) });
  const games = useQuery({ queryKey: ["user", id, "games"], queryFn: () => usersApi.games(id, 100, 0) });
  const series = useQuery({
    queryKey: ["user", id, "series", "preview"],
    queryFn: () => usersApi.series(id, { limit: 3, offset: 0, show_past: true, show_closed: true }),
  });
  const club = useQuery({
    queryKey: ["club", user.data?.club_id],
    queryFn: () => clubsApi.get(user.data!.club_id!),
    enabled: !!user.data?.club_id,
  });

  if (user.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (user.error) return <PageShell><ErrorBlock error={user.error} /></PageShell>;
  if (!user.data) return null;
  if (location.pathname !== `/user/${id}`) {
    return <Outlet />;
  }

  const allGames = games.data?.items ?? [];
  const allSeries = series.data?.items ?? [];
  const recentGames = allGames.slice(0, 3);
  const visibleSeries = allSeries;
  const canBlockFromProfile =
    !!me?.club_id &&
    (me.club_state === ClubState.Leader || me.club_state === ClubState.President) &&
    me.id !== user.data.id;

  const blockFromProfile = async () => {
    if (!me?.club_id) return;
    try {
      await clubsApi.blockProfile(me.club_id, user.data.id);
      qc.invalidateQueries({ queryKey: ["club", me.club_id, "members"] });
      qc.invalidateQueries({ queryKey: ["club", me.club_id, "bans"] });
      toast.success("Игрок заблокирован");
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    }
  };

  return (
    <PageShell>
      <PageHeader
        eyebrow="Игрок"
        title={displayUserName(user.data)}
        description={user.data.description}
        actions={
          canBlockFromProfile ? (
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="secondary">Заблокировать в своем клубе</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Заблокировать игрока?</AlertDialogTitle>
                  <AlertDialogDescription>
                    Игрок будет заблокирован в вашем клубе и не сможет вступить в него.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Отмена</AlertDialogCancel>
                  <AlertDialogAction onClick={() => void blockFromProfile()}>
                    Заблокировать
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          ) : undefined
        }
      />

      <div className="grid gap-6 lg:grid-cols-3">
        <aside className="space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-3 font-display text-lg font-semibold">Профиль</h2>
            {user.data.name && user.data.show_name && (
              <p className="text-sm"><span className="text-muted-foreground">Имя: </span>{user.data.name}</p>
            )}
            {user.data.nickname && (
              <p className="text-sm"><span className="text-muted-foreground">Никнейм: </span>{user.data.nickname}</p>
            )}
          </div>
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-3 font-display text-lg font-semibold">Клуб</h2>
            {user.data.club_id ? (
              <div className="space-y-2">
                <Link to="/clubs/$id" params={{ id: user.data.club_id }} className="block font-semibold text-primary hover:underline">
                  {club.data?.name ?? "Открыть клуб"}
                </Link>
                <RoleBadge state={(user.data.club_state ?? ClubState.None) as ClubState} />
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Без клуба</p>
            )}
          </div>
        </aside>

        <section className="space-y-6 lg:col-span-2">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="font-display text-xl font-semibold">Игры</h2>
              <Button size="sm" variant="outline" asChild>
                <Link to="/user/$id/games" params={{ id }}>Посмотреть все игры игрока</Link>
              </Button>
            </div>
            {!recentGames.length ? <EmptyBlock title="Игр нет" /> : (
              <div className="space-y-2">
                {recentGames.map((g) => (
                  <div key={g.id} className="flex items-center justify-between gap-3 rounded-lg border border-border/40 bg-background/40 p-3">
                    <Link to="/game/$id" params={{ id: g.id }} className="min-w-0 flex-1 hover:text-primary">
                      <p className="font-medium">{g.name || `Игра #${g.number}`}</p>
                      <p className="text-xs text-muted-foreground">{g.series_name}</p>
                      <p className="text-xs text-muted-foreground">{fmtDate(g.created_at)}</p>
                    </Link>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="font-display text-xl font-semibold">Серии</h2>
              <Button size="sm" variant="outline" asChild>
                <Link to="/user/$id/series" params={{ id }}>Посмотреть все серии игрока</Link>
              </Button>
            </div>
            {!visibleSeries.length ? <EmptyBlock title="Серий нет" /> : (
              <div className="space-y-2">
                {visibleSeries.map((s) => (
                  <div key={s.id} className="flex items-center justify-between gap-3 rounded-lg border border-border/40 bg-background/40 p-3">
                    <Link to="/series/$id" params={{ id: s.id }} className="min-w-0 flex-1 hover:text-primary">
                      <p className="font-medium">{s.name}</p>
                      <p className="text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
                      <div className="mt-1 flex flex-wrap items-center gap-2">
                        {s.is_rating && (
                          <span className="inline-flex rounded-full bg-sky-100 px-2 py-0.5 text-xs text-sky-800">
                            На рейтинг
                          </span>
                        )}
                        {Number(s.price_rub ?? 0) > 0 && (
                          <span className="inline-flex rounded-full bg-amber-100 px-2 py-0.5 text-xs text-amber-800">
                            Платно · {fmtRub(s.price_rub)}
                          </span>
                        )}
                      </div>
                    </Link>
                  </div>
                ))}
              </div>
            )}
          </div>
        </section>
      </div>
    </PageShell>
  );
}
