import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { clubsApi, authApi, ApiError, seriesApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { RoleBadge } from "@/components/site/RoleBadge";
import { ClubState } from "@/types/api";
import { displayUserName, canManageClub } from "@/lib/roles";
import { useAuthStore } from "@/lib/auth-store";
import { Button } from "@/components/ui/button";
import { fmtDateRange } from "@/lib/format";
import { Settings, UserPlus } from "lucide-react";
import { toast } from "sonner";
import { useMemo } from "react";

export const Route = createFileRoute("/clubs/$id")({ component: ClubPage });

function ClubPage() {
  const { id } = Route.useParams();
  const location = useLocation();
  const me = useAuthStore((s) => s.me);
  const setMe = useAuthStore((s) => s.setMe);
  const qc = useQueryClient();

  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const members = useQuery({ queryKey: ["club", id, "members"], queryFn: () => clubsApi.members(id) });
  const series = useQuery({ queryKey: ["club", id, "series"], queryFn: () => clubsApi.series(id) });
  const seriesIds = useMemo(() => (series.data?.items ?? []).map((s) => s.id), [series.data]);
  const games = useQuery({
    queryKey: ["club", id, "games", seriesIds],
    enabled: seriesIds.length > 0,
    queryFn: async () => {
      const responses = await Promise.all(seriesIds.map((sid) => seriesApi.games(sid).catch(() => null)));
      return responses
        .flatMap((res, i) => (res?.items ?? []).map((g) => ({ ...g, seriesName: series.data!.items![i].name })))
        .sort((a, b) => b.number - a.number);
    },
  });

  const canManage = canManageClub(me, id);
  const isInClub = me?.club_id === id;
  const canJoin = me && !me.club_id;

  const onJoin = async () => {
    try {
      await clubsApi.join(id);
      const u = await authApi.me();
      setMe(u);
      qc.invalidateQueries({ queryKey: ["club", id, "members"] });
      toast.success("Вы вступили в клуб");
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Не удалось вступить в клуб");
    }
  };
  const deleteSeries = async (seriesId: string) => {
    if (!confirm("Удалить эту серию?")) return;
    try {
      await (await import("@/lib/api")).seriesApi.delete(seriesId);
      qc.invalidateQueries({ queryKey: ["club", id, "series"] });
      toast.success("Серия удалена");
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    }
  };

  if (club.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (club.error) return <PageShell><ErrorBlock error={club.error} /></PageShell>;
  if (!club.data) return null;
  if (location.pathname !== `/clubs/${id}`) {
    return <Outlet />;
  }

  return (
    <PageShell>
      <PageHeader
        eyebrow="Клуб" title={club.data.name} description={club.data.description}
        actions={
          <>
            {me && !isInClub && me.club_id && (
              <span className="rounded-md bg-muted px-3 py-2 text-xs text-muted-foreground">Вы уже состоите в другом клубе</span>
            )}
            {canJoin && (
              <Button onClick={onJoin}><UserPlus className="mr-1 h-4 w-4" />Вступить в клуб</Button>
            )}
            {canManage && (
              <Button asChild variant="outline"><Link to="/clubs/$id/manage" params={{ id }}><Settings className="mr-1 h-4 w-4" />Управление</Link></Button>
            )}
          </>
        }
      />

      <div className="grid gap-6 lg:grid-cols-3">
        <section className="lg:col-span-2 space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="font-display text-xl font-semibold">Серии</h2>
              <div className="flex items-center gap-2">
                <Button size="sm" variant="outline" asChild><Link to="/clubs/$id/series" params={{ id }}>Все серии</Link></Button>
                {canManage && (
                  <Button size="sm" asChild><Link to="/series/create">Создать серию</Link></Button>
                )}
              </div>
            </div>
            {series.isLoading ? <LoadingBlock /> :
              !series.data?.items?.length ? <EmptyBlock title="Серий пока нет" /> : (
              <div className="space-y-2">
                {series.data.items.map((s) => (
                  <div key={s.id} className="flex items-center justify-between gap-3 rounded-lg border border-border/40 bg-background/40 p-3">
                    <Link to="/series/$id" params={{ id: s.id }} className="min-w-0 flex-1 hover:text-primary">
                      <p className="font-medium">{s.name}</p>
                      <p className="text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
                    </Link>
                    {canManage && (
                      <Button size="sm" variant="outline" onClick={() => void deleteSeries(s.id)}>
                        Удалить
                      </Button>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="font-display text-xl font-semibold">Текущие игры</h2>
              <Button size="sm" variant="outline" asChild><Link to="/clubs/$id/games" params={{ id }}>Все игры</Link></Button>
            </div>
            {games.isLoading ? <LoadingBlock /> :
              !games.data?.length ? <EmptyBlock title="Игр пока нет" /> : (
              <div className="space-y-2">
                {games.data.slice(0, 5).map((g) => (
                  <Link key={g.id} to="/game/$id" params={{ id: g.id }}
                    className="flex items-center justify-between rounded-lg border border-border/40 bg-background/40 p-3 hover:border-primary/50">
                    <div>
                      <p className="font-medium">{g.name || `Игра #${g.number}`}</p>
                      <p className="text-xs text-muted-foreground">{g.seriesName}</p>
                    </div>
                    <span className="text-xs text-muted-foreground">#{g.number}</span>
                  </Link>
                ))}
              </div>
            )}
          </div>
        </section>

        <aside className="space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-4 font-display text-lg font-semibold">Участники</h2>
            {members.isLoading ? <LoadingBlock /> :
              !members.data?.items?.length ? <p className="text-sm text-muted-foreground">Участников пока нет</p> : (
              <ul className="space-y-2">
                {members.data.items.map((m) => (
                  <li key={m.id} className="flex items-center justify-between gap-2 text-sm">
                    <Link to="/user/$id" params={{ id: m.id }} className="truncate hover:text-primary">
                      {displayUserName(m)}
                    </Link>
                    <RoleBadge state={(m.club_state ?? ClubState.None) as ClubState} />
                  </li>
                ))}
              </ul>
            )}
          </div>
        </aside>
      </div>
    </PageShell>
  );
}
