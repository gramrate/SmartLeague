import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { gamesApi, seriesApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock } from "@/components/site/States";
import { displayUserName, canManageClub } from "@/lib/roles";
import { useAuthStore } from "@/lib/auth-store";
import { Button } from "@/components/ui/button";
import { Settings } from "lucide-react";
import { useMemo } from "react";

export const Route = createFileRoute("/game/$id")({ component: GamePage });

const ROLE_LABEL: Record<string, string> = {
  civilian: "Мирный",
  mafia: "Мафия",
  don: "Дон",
  sheriff: "Шериф",
};

function GamePage() {
  const { id } = Route.useParams();
  const location = useLocation();
  const me = useAuthStore((s) => s.me);
  const game = useQuery({ queryKey: ["game", id, "full"], queryFn: () => gamesApi.full(id) });
  const series = useQuery({
    queryKey: ["series", game.data?.series_id],
    queryFn: () => seriesApi.get(game.data!.series_id),
    enabled: !!game.data?.series_id,
  });
  const participants = useQuery({
    queryKey: ["series", game.data?.series_id, "participants"],
    queryFn: () => seriesApi.participants(game.data!.series_id),
    enabled: !!game.data?.series_id,
  });

  const canManage = !!series.data && canManageClub(me, series.data.club_id);

  const byId = useMemo(
    () => new Map((participants.data?.items ?? []).map((u) => [u.id, u])),
    [participants.data],
  );

  if (game.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (game.error) return <PageShell><ErrorBlock error={game.error} /></PageShell>;
  if (!game.data) return null;
  if (location.pathname !== `/game/${id}`) {
    return <Outlet />;
  }

  const results = (game.data.results ?? []).slice().sort((a, b) => (a.place ?? 99) - (b.place ?? 99));
  return (
    <PageShell>
      <PageHeader
        eyebrow={series.data?.name ?? "Игра"}
        title={game.data.name || `Игра #${game.data.number}`}
        description={game.data.description}
        actions={
          <>
            {series.data && (
              <Button variant="outline" asChild>
                <Link to="/series/$id" params={{ id: series.data.id }}>Открыть серию</Link>
              </Button>
            )}
            {canManage && (
              <Button asChild>
                <Link to="/game/$id/manage" params={{ id }}><Settings className="mr-1 h-4 w-4" />Управление</Link>
              </Button>
            )}
          </>
        }
      />

      <div className="grid gap-6">
        <section className="w-full rounded-2xl border border-border/60 bg-card/60 p-6">
          <h2 className="mb-4 font-display text-xl font-semibold">Результаты</h2>
          {results.length === 0 ? (
            <p className="text-sm text-muted-foreground">Результатов пока нет.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[760px] text-sm">
                <thead>
                  <tr className="border-b border-border/60 text-left text-xs uppercase text-muted-foreground">
                    <th className="py-2 pr-2">Место</th>
                    <th className="py-2 pr-2">Никнейм</th>
                    <th className="py-2 pr-2">Роль</th>
                    <th className="py-2 pr-2">Лучший ход</th>
                    <th className="py-2 pr-2">Компенсация</th>
                    <th className="py-2 pr-2">ЖК</th>
                    <th className="py-2 pr-2">Удаление</th>
                    <th className="py-2 pr-2">Доп балл</th>
                    <th className="py-2 pr-2 text-right">Итоговый балл</th>
                  </tr>
                </thead>
                <tbody>
                  {results.map((r) => (
                    <tr key={r.profile_id} className="border-b border-border/30">
                      <td className="py-2 pr-2 font-bold text-primary">{r.place ?? "—"}</td>
                      <td className="py-2 pr-2">
                        <Link to="/user/$id" params={{ id: r.profile_id }} className="inline-block max-w-[180px] truncate align-bottom hover:text-primary">
                          {displayUserName(byId.get(r.profile_id) ?? { id: r.profile_id })}
                        </Link>
                      </td>
                      <td className="py-2 pr-2 text-muted-foreground">{r.role ? (ROLE_LABEL[r.role] ?? r.role) : "—"}</td>
                      <td className="py-2 pr-2">{r.best_move ?? "—"}</td>
                      <td className="py-2 pr-2">{r.compensation?.toFixed(2) ?? "—"}</td>
                      <td className="py-2 pr-2">{r.yellow_cards?.toFixed(2) ?? "—"}</td>
                      <td className="py-2 pr-2">{r.removed?.toFixed(2) ?? "—"}</td>
                      <td className="py-2 pr-2">{r.extra_points?.toFixed(2) ?? "—"}</td>
                      <td className="py-2 pr-2 text-right font-mono">{r.total_points?.toFixed(2) ?? "—"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </PageShell>
  );
}
