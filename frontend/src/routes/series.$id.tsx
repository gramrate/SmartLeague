import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { seriesApi, clubsApi, ApiError } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { fmtDateRange, fmtRub } from "@/lib/format";
import { useAuthStore } from "@/lib/auth-store";
import { canManageClub, displayUserName } from "@/lib/roles";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";
import { toast } from "sonner";
import { UserLink } from "@/components/site/UserLink";
import { ClubState, GameStatus } from "@/types/api";
import { RoleBadge } from "@/components/site/RoleBadge";

export const Route = createFileRoute("/series/$id")({ component: SeriesPage });

function SeriesPage() {
  const { id } = Route.useParams();
  const me = useAuthStore((s) => s.me);
  const qc = useQueryClient();
  const full = useQuery({ queryKey: ["series", id, "full"], queryFn: () => seriesApi.full(id) });

  const series = full.data?.series;
  const club = useQuery({
    queryKey: ["club", series?.club_id],
    queryFn: () => clubsApi.get(series!.club_id),
    enabled: !!series?.club_id,
  });

  const canManage = !!series && canManageClub(me, series.club_id);
  const isParticipant = !!me && (full.data?.participants.items ?? []).some((p) => p.id === me.id);

  const [gameName, setGameName] = useState("");
  const [gameDescription, setGameDescription] = useState("");
  const [creatingGame, setCreatingGame] = useState(false);
  const [joining, setJoining] = useState(false);
  const [leaving, setLeaving] = useState(false);
  const [updatingSeriesStatus, setUpdatingSeriesStatus] = useState(false);
  const createGame = async () => {
    setCreatingGame(true);
    try {
      const payload = {
        name: gameName.trim() || undefined,
        description: gameDescription.trim() || undefined,
        status: GameStatus.Draft,
      };
      await seriesApi.createGameDraft(id, payload);
      qc.invalidateQueries({ queryKey: ["series", id, "full"] });
      setGameName("");
      setGameDescription("");
      toast.success("Игра создана как черновик");
    } catch (err) { toast.error(err instanceof ApiError ? err.message : "Ошибка"); }
    finally { setCreatingGame(false); }
  };
  const joinSeries = async () => {
    setJoining(true);
    try {
      await seriesApi.join(id);
      qc.invalidateQueries({ queryKey: ["series", id, "full"] });
      toast.success("Вы присоединились к серии");
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setJoining(false);
    }
  };
  const leaveSeries = async () => {
    setLeaving(true);
    try {
      await seriesApi.leave(id);
      qc.invalidateQueries({ queryKey: ["series", id, "full"] });
      toast.success("Вы вышли из серии");
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setLeaving(false);
    }
  };
  const toggleSeriesRegistration = async () => {
    if (!series) return;
    const nextClosed = !series.is_closed;
    if (!confirm(nextClosed ? "Закрыть регистрацию в серии?" : "Открыть регистрацию в серии?")) return;
    setUpdatingSeriesStatus(true);
    try {
      await seriesApi.update(series.id, { is_closed: nextClosed });
      qc.invalidateQueries({ queryKey: ["series", id, "full"] });
      qc.invalidateQueries({ queryKey: ["series", id] });
      toast.success(nextClosed ? "Серия закрыта" : "Серия открыта");
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setUpdatingSeriesStatus(false);
    }
  };

  if (full.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (full.error) return <PageShell><ErrorBlock error={full.error} /></PageShell>;
  if (!series) return null;

  const participantsById = new Map((full.data!.participants.items ?? []).map((p) => [p.id, p]));

  return (
    <PageShell>
      <PageHeader
        eyebrow={club.data ? club.data.name : "Серия"}
        title={series.name}
        description={series.description}
        actions={
          <>
            {club.data && (
              <Button variant="outline" asChild><Link to="/clubs/$id" params={{ id: series.club_id }}>{club.data.name}</Link></Button>
            )}
          </>
        }
      />
      <div className="mb-3 flex items-center gap-2 text-sm">
        <span className="text-muted-foreground">Статус регистрации:</span>
        <span className={series.is_closed ? "text-destructive" : "text-emerald-600"}>
          {series.is_closed ? "Закрыта" : "Открыта"}
        </span>
      </div>
      <div className="mb-6 flex flex-wrap items-center gap-2 rounded-xl border border-border/60 bg-card/40 p-3">
        {me && !isParticipant && (
          <Button onClick={joinSeries} disabled={joining || !!series.is_closed}>
            {joining ? "Подключение..." : "Присоединиться к серии"}
          </Button>
        )}
        {me && isParticipant && (
          <Button variant="outline" onClick={leaveSeries} disabled={leaving || !!series.is_closed}>
            {leaving ? "Выход..." : "Покинуть серию"}
          </Button>
        )}
        {canManage && (
          <Button className="ml-auto" variant="outline" disabled={updatingSeriesStatus} onClick={toggleSeriesRegistration}>
            {updatingSeriesStatus ? "Сохранение..." : series.is_closed ? "Открыть регистрацию" : "Закрыть регистрацию"}
          </Button>
        )}
      </div>
      <p className="mb-6 text-sm text-muted-foreground">{fmtDateRange(series.start_at, series.end_at)}</p>
      <div className="mb-6 flex flex-wrap items-center gap-2">
        {Number(series.price_rub ?? 0) > 0 && (
          <div className="inline-flex rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800">
            Платно · {fmtRub(series.price_rub)}
          </div>
        )}
        <div className="inline-flex rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-800">
          {series.is_rating ? "На рейтинг" : "Без рейтинга"}
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <section className="lg:col-span-2 space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-4 font-display text-xl font-semibold">Игры</h2>
            {!full.data!.games.items?.length ? <EmptyBlock title="Игр пока нет" /> : (
              <ul className="divide-y divide-border/40">
                {full.data!.games.items.map((g) => (
                  <li key={g.id}>
                    <Link to="/game/$id" params={{ id: g.id }} className="flex items-center justify-between py-3 hover:text-primary">
                      <span>#{g.number} · {g.name || "Игра"}</span>
                      <span className="text-xs text-muted-foreground">{g.status === 2 ? "Завершена" : g.status === 1 ? "Идет" : "Черновик"}</span>
                    </Link>
                  </li>
                ))}
              </ul>
            )}
            {canManage && (
              <form className="mt-6 border-t border-border/40 pt-4">
                <Label className="text-xs">Новая игра</Label>
                <div className="mt-2 space-y-3">
                  <Input placeholder="Название игры" value={gameName} onChange={(e) => setGameName(e.target.value)} />
                  <Textarea
                    rows={3}
                    placeholder="Описание игры"
                    value={gameDescription}
                    onChange={(e) => setGameDescription(e.target.value)}
                  />
                  <Button type="button" disabled={creatingGame} onClick={() => void createGame()}>
                    {creatingGame ? "Создание..." : "Создать"}
                  </Button>
                </div>
              </form>
            )}
          </div>

          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-4 font-display text-xl font-semibold">Таблица лидеров</h2>
            {!full.data!.leaderboard.items?.length ? <p className="text-sm text-muted-foreground">Пусто</p> : (
              <ol className="divide-y divide-border/40">
                {full.data!.leaderboard.items.map((r, i) => (
                  <li key={r.profile_id} className="flex items-center justify-between py-3">
                    <span className="flex items-center gap-3">
                      <span className="grid h-7 w-7 place-items-center rounded-full bg-primary/15 text-xs font-bold text-primary">{i + 1}</span>
                      <Link to="/user/$id" params={{ id: r.profile_id }} className="hover:text-primary">
                        {displayUserName(participantsById.get(r.profile_id) ?? { id: r.profile_id })}
                      </Link>
                    </span>
                    <span className="font-mono text-sm">{r.points.toFixed(2)}</span>
                  </li>
                ))}
              </ol>
            )}
          </div>
        </section>

        <aside className="space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-4 font-display text-lg font-semibold">Участники</h2>
            {!full.data!.participants.items?.length ? <p className="text-sm text-muted-foreground">Участников нет</p> : (
              <ul className="space-y-2 text-sm">
                {full.data!.participants.items.map((p) => (
                  <li key={p.id} className="flex items-center justify-between gap-2">
                    <Link to="/user/$id" params={{ id: p.id }} className="truncate hover:text-primary">{displayUserName(p)}</Link>
                    {p.club_id === series.club_id && (p.club_state ?? ClubState.None) !== ClubState.None ? (
                      <RoleBadge state={(p.club_state ?? ClubState.None) as ClubState} />
                    ) : (
                      <span />
                    )}
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
