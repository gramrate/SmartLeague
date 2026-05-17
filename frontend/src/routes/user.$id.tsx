import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { usersApi, clubsApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { RoleBadge } from "@/components/site/RoleBadge";
import { ClubState } from "@/types/api";
import { displayUserName } from "@/lib/roles";
import { fmtDate, fmtDateRange } from "@/lib/format";
import { useState } from "react";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/user/$id")({ component: UserPage });

function UserPage() {
  const { id } = Route.useParams();
  const user = useQuery({ queryKey: ["user", id], queryFn: () => usersApi.get(id) });
  const games = useQuery({ queryKey: ["user", id, "games"], queryFn: () => usersApi.games(id, 100, 0) });
  const series = useQuery({ queryKey: ["user", id, "series"], queryFn: () => usersApi.series(id, 100, 0) });
  const club = useQuery({
    queryKey: ["club", user.data?.club_id],
    queryFn: () => clubsApi.get(user.data!.club_id!),
    enabled: !!user.data?.club_id,
  });

  const [showAllGames, setShowAllGames] = useState(false);
  const [showAllSeries, setShowAllSeries] = useState(false);

  if (user.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (user.error) return <PageShell><ErrorBlock error={user.error} /></PageShell>;
  if (!user.data) return null;

  const allGames = games.data?.items ?? [];
  const allSeries = series.data?.items ?? [];
  const recentGames = showAllGames ? allGames : allGames.slice(0, 3);
  const visibleSeries = showAllSeries ? allSeries : allSeries.slice(0, 5);

  return (
    <PageShell>
      <PageHeader eyebrow="Игрок" title={displayUserName(user.data)} description={user.data.description} />

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
              <h2 className="font-display text-xl font-semibold">Последние игры</h2>
              {allGames.length > 3 && (
                <Button size="sm" variant="ghost" onClick={() => setShowAllGames((v) => !v)}>
                  {showAllGames ? "Показать последние" : `Показать все (${allGames.length})`}
                </Button>
              )}
            </div>
            {!recentGames.length ? <EmptyBlock title="Игр нет" /> : (
              <ul className="divide-y divide-border/40">
                {recentGames.map((g) => (
                  <li key={g.id}>
                    <Link to="/game/$id" params={{ id: g.id }} className="flex items-center justify-between py-3 hover:text-primary">
                      <span>
                        <span className="font-medium">{g.name || `Игра #${g.number}`}</span>
                        <span className="ml-2 text-xs text-muted-foreground">{g.series_name}</span>
                      </span>
                      <span className="text-xs text-muted-foreground">{fmtDate(g.created_at)}</span>
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="font-display text-xl font-semibold">Серии</h2>
              {allSeries.length > 5 && (
                <Button size="sm" variant="ghost" onClick={() => setShowAllSeries((v) => !v)}>
                  {showAllSeries ? "Свернуть" : `Показать все (${allSeries.length})`}
                </Button>
              )}
            </div>
            {!visibleSeries.length ? <EmptyBlock title="Серий нет" /> : (
              <ul className="divide-y divide-border/40">
                {visibleSeries.map((s) => (
                  <li key={s.id}>
                    <Link to="/series/$id" params={{ id: s.id }} className="flex items-center justify-between py-3 hover:text-primary">
                      <span className="font-medium">{s.name}</span>
                      <span className="text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</span>
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </section>
      </div>
    </PageShell>
  );
}
