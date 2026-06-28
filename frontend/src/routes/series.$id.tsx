import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { seriesApi, clubsApi, gamesApi, ApiError } from "@/lib/api";
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
import { ClubState, GameStatus, JudgeRole } from "@/types/api";
import { RoleBadge } from "@/components/site/RoleBadge";
import { Settings } from "lucide-react";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { X } from "lucide-react";

export const Route = createFileRoute("/series/$id")({ component: SeriesPage });

function SeriesPage() {
  const { id } = Route.useParams();
  const location = useLocation();
  const navigate = useNavigate();
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
  const isBanned = full.data?.is_banned ?? false;
  const judges = full.data?.judges ?? [];
  const myJudge = me ? judges.find((j) => j.profile_id === me.id) : undefined;
  const isJudge = !!myJudge;

  const judgesQ = useQuery({
    queryKey: ["series", id, "judges"],
    queryFn: () => seriesApi.judges(id),
    enabled: canManage,
  });
  const allParticipants = useQuery({
    queryKey: ["series", id, "participants", "manage-payments"],
    queryFn: () => seriesApi.participants(id, 200, 0),
    enabled: canManage,
  });

  const [judgeDialogOpen, setJudgeDialogOpen] = useState(false);
  const [judgeSearchQ, setJudgeSearchQ] = useState("");
  const [settingJudge, setSettingJudge] = useState(false);
  const [removingJudge, setRemovingJudge] = useState<string | null>(null);

  const invalidateJudgeRelated = () => {
    qc.invalidateQueries({ queryKey: ["series", id, "judges"] });
    qc.invalidateQueries({ queryKey: ["series", id, "full"] });
    qc.invalidateQueries({ queryKey: ["series", id, "participants"], exact: true });
    qc.invalidateQueries({ queryKey: ["series", id, "participants", "manage-payments"], exact: true });
    qc.invalidateQueries({ queryKey: ["series", id, "payments"], exact: true });
  };
  const judgeSet = (profileId: string, role: JudgeRole) => async () => {
    setSettingJudge(true);
    try {
      await seriesApi.setJudge(id, profileId, role);
      invalidateJudgeRelated();
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
    finally { setSettingJudge(false); }
  };
  const judgeRemove = (profileId: string) => async () => {
    setRemovingJudge(profileId);
    try {
      await seriesApi.removeJudge(id, profileId);
      invalidateJudgeRelated();
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
    finally { setRemovingJudge(null); }
  };

  const [gameName, setGameName] = useState("");
  const [gameDescription, setGameDescription] = useState("");
  const gameNameLimit = 100;
  const gameDescriptionLimit = 500;
  const [creatingGame, setCreatingGame] = useState(false);
  const [joining, setJoining] = useState(false);
  const [leaving, setLeaving] = useState(false);
  const [deletingGameId, setDeletingGameId] = useState<string | null>(null);
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
  const deleteGame = async (gameId: string) => {
    setDeletingGameId(gameId);
    try {
      await gamesApi.delete(gameId);
      qc.invalidateQueries({ queryKey: ["series", id, "full"] });
      toast.success("Игра удалена");
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Ошибка");
    } finally {
      setDeletingGameId(null);
    }
  };
  if (full.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (full.error) return <PageShell><ErrorBlock error={full.error} /></PageShell>;
  if (!series) return null;
  if (location.pathname !== `/series/${id}`) {
    return <Outlet />;
  }

  const participantsById = new Map((full.data!.participants.items ?? []).map((p) => [p.id, p]));

  return (
    <PageShell>
      <PageHeader
        eyebrow={
          club.data
            ? <Link to="/clubs/$id" params={{ id: series.club_id }} className="hover:underline">{club.data.name}</Link>
            : "Серия"
        }
        title={series.name}
        description={series.description}
        actions={
          <>
            {club.data && (
              <Button variant="outline" onClick={() => {
                qc.invalidateQueries({ queryKey: ["club", series.club_id] });
                void navigate({ to: "/clubs/$id", params: { id: series.club_id } });
              }}>К клубу</Button>
            )}
            {canManage && (
              <Button asChild>
                <Link to="/series/$id/manage" params={{ id: series.id }}>
                  <Settings className="mr-1 h-4 w-4" />
                  Управление серией
                </Link>
              </Button>
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
        <span className="text-sm text-muted-foreground">Участников: {full.data?.participants.items?.length ?? 0}</span>
      </div>
      <p className="mb-6 text-sm text-muted-foreground">{fmtDateRange(series.start_at, series.end_at)}</p>
      <div className="mb-6 flex flex-wrap items-center gap-2">
        {series.is_tournament && (
          <div className="inline-flex rounded-full bg-purple-100 px-3 py-1 text-xs font-medium text-purple-800 dark:bg-purple-900/40 dark:text-purple-300">
            Турнир
          </div>
        )}
        {series.is_club_only && (
          <div className="inline-flex rounded-full bg-teal-100 px-3 py-1 text-xs font-medium text-teal-800 dark:bg-teal-900/40 dark:text-teal-300">
            Для участников клуба
          </div>
        )}
        {Number(series.price_rub ?? 0) > 0 && (
          <div className="inline-flex rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800">
            Платно · {fmtRub(series.price_rub)}
          </div>
        )}
        {!series.is_tournament && (
          <div className="inline-flex rounded-full bg-sky-100 px-3 py-1 text-xs font-medium text-sky-800">
            {series.is_rating ? "На рейтинг" : "Без рейтинга"}
          </div>
        )}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <section className="lg:col-span-2 space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-4 font-display text-xl font-semibold">Игры</h2>
            {!full.data!.games.items?.length ? <EmptyBlock title="Игр пока нет" /> : (
              <ul className="divide-y divide-border/40">
                {full.data!.games.items.map((g) => (
                  <li key={g.id}>
                    <div className="flex items-center justify-between gap-3 py-3">
                      <Link to="/game/$id" params={{ id: g.id }} className="min-w-0 flex-1 hover:text-primary">
                        <div className="flex items-center justify-between gap-2">
                          <span className="truncate">#{g.number} · {g.name || "Игра"}</span>
                          <span className="shrink-0 text-xs text-muted-foreground">{g.status === 2 ? "Завершена" : g.status === 1 ? "Идет" : "Черновик"}</span>
                        </div>
                      </Link>
                      {canManage && (
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button variant="destructive" size="sm">Удалить</Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Удалить игру?</AlertDialogTitle>
                              <AlertDialogDescription>
                                Это действие нельзя отменить. Игра будет удалена из серии.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Отмена</AlertDialogCancel>
                              <AlertDialogAction
                                variant="destructive"
                                onClick={() => void deleteGame(g.id)}
                                disabled={deletingGameId === g.id}
                              >
                                {deletingGameId === g.id ? "Удаление..." : "Удалить"}
                              </AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      )}
                    </div>
                  </li>
                ))}
              </ul>
            )}
            {canManage && (
              <form className="mt-6 border-t border-border/40 pt-4">
                <Label className="text-xs">Новая игра</Label>
                <div className="mt-2 space-y-3">
                  <div className="space-y-1">
                    <Input placeholder="Название игры" value={gameName} onChange={(e) => setGameName(e.target.value)} maxLength={gameNameLimit} />
                    <p className="text-xs text-muted-foreground">{gameName.length}/{gameNameLimit}</p>
                  </div>
                  <Textarea
                    rows={3}
                    placeholder="Описание игры"
                    value={gameDescription}
                    onChange={(e) => setGameDescription(e.target.value)}
                    maxLength={gameDescriptionLimit}
                  />
                  <p className="text-xs text-muted-foreground">{gameDescription.length}/{gameDescriptionLimit}</p>
                  <Button type="button" disabled={creatingGame || gameName.length > gameNameLimit || gameDescription.length > gameDescriptionLimit} onClick={() => void createGame()}>
                    {creatingGame ? "Создание..." : "Создать"}
                  </Button>
                </div>
              </form>
            )}
          </div>

          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="font-display text-xl font-semibold">Таблица лидеров</h2>
              {(full.data!.leaderboard.pagination?.total_items ?? 0) > 10 && (
                <Link to="/series/$id/leaderboard" params={{ id }} className="text-sm text-primary hover:underline">
                  Все ({full.data!.leaderboard.pagination.total_items})
                </Link>
              )}
            </div>
            {!full.data!.leaderboard.items?.length ? <p className="text-sm text-muted-foreground">Пусто</p> : (
              <ol className="divide-y divide-border/40">
                {full.data!.leaderboard.items.map((r, i) => (
                  <li key={r.profile_id ?? r.guest_id ?? i} className="flex items-center justify-between py-3">
                    <span className="flex items-center gap-3">
                      <span className="grid h-7 w-7 place-items-center rounded-full bg-primary/15 text-xs font-bold text-primary">{i + 1}</span>
                      {r.profile_id ? (
                        <Link to="/user/$id" params={{ id: r.profile_id }} className="max-w-[220px] truncate hover:text-primary">
                          {displayUserName(participantsById.get(r.profile_id) ?? { id: r.profile_id })}
                        </Link>
                      ) : (
                        <span className="max-w-[220px] truncate text-muted-foreground" title="Гость">
                          {r.guest_nickname ?? "Гость"}
                        </span>
                      )}
                    </span>
                    <span className="font-mono text-sm">{r.points.toFixed(2)}</span>
                  </li>
                ))}
              </ol>
            )}
          </div>
        </section>

        <aside className="space-y-6">
          {(canManage || judges.length > 0) && (
            <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
              <div className="mb-4 flex items-center justify-between gap-2">
                <h2 className="font-display text-lg font-semibold">Судьи</h2>
                {canManage && (
                  <Button size="sm" variant="outline" onClick={() => setJudgeDialogOpen(true)}>
                    Изменить состав
                  </Button>
                )}
              </div>
              {judges.length === 0 ? (
                <p className="text-sm text-muted-foreground">Судей пока нет.</p>
              ) : (
                <ul className="space-y-2 text-sm">
                  {judges.map((j) => (
                    <li key={j.profile_id} className="flex items-center justify-between gap-2">
                      <Link to="/user/$id" params={{ id: j.profile_id }} className="min-w-0 flex-1 truncate hover:text-primary">
                        {displayUserName(j)}
                      </Link>
                      <span className={`shrink-0 rounded-full px-2 py-0.5 text-xs font-medium ${j.role === JudgeRole.Main ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-300" : "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300"}`}>
                        {j.role === JudgeRole.Main ? "Главный" : "Боковой"}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          )}
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            {me && (
              <div className="mb-4 flex flex-col items-end gap-2">
                {isJudge || isParticipant ? (
                  <Button variant="outline" onClick={leaveSeries} disabled={leaving || (isParticipant && !!series.is_closed)}>
                    {leaving ? "Выход..." : "Покинуть серию"}
                  </Button>
                ) : isBanned ? (
                  <>
                    <Button disabled>Присоединиться к серии</Button>
                    <p className="text-xs text-destructive">Вы заблокированы в этом клубе</p>
                  </>
                ) : !series.is_closed ? (
                  <Button onClick={joinSeries} disabled={joining}>
                    {joining ? "Подключение..." : "Присоединиться к серии"}
                  </Button>
                ) : null}
              </div>
            )}
            <h2 className="mb-4 font-display text-lg font-semibold">Участники</h2>
            {!full.data!.participants.items?.length ? <p className="text-sm text-muted-foreground">Участников нет</p> : (
              <ul className="space-y-2 text-sm">
                {full.data!.participants.items.map((p) => (
                  <li key={p.id} className="flex items-center justify-between gap-2">
                    <Link to="/user/$id" params={{ id: p.id }} className="min-w-0 flex-1 truncate hover:text-primary">{displayUserName(p)}</Link>
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

      <Dialog open={judgeDialogOpen} onOpenChange={(open) => { setJudgeDialogOpen(open); if (!open) setJudgeSearchQ(""); }}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Управление судьями</DialogTitle>
            <DialogDescription>Назначайте судей из участников серии.</DialogDescription>
          </DialogHeader>
          <Tabs defaultValue="judges" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="judges">Текущие судьи</TabsTrigger>
              <TabsTrigger value="participants">Участники</TabsTrigger>
            </TabsList>
            <TabsContent value="judges" className="mt-3">
              <div className="rounded-lg border border-border/40 bg-background/40 p-3">
                {!judgesQ.data?.items?.length ? (
                  <p className="text-sm text-muted-foreground">Судей пока нет.</p>
                ) : (
                  <ul className="divide-y divide-border/40">
                    {judgesQ.data.items.map((j) => (
                      <li key={j.profile_id} className="flex items-center justify-between gap-2 py-2 text-sm">
                        <div className="flex min-w-0 items-center gap-2">
                          <Link to="/user/$id" params={{ id: j.profile_id }} className="truncate hover:text-primary">{displayUserName(j)}</Link>
                          <span className={`shrink-0 rounded-full px-2 py-0.5 text-xs font-medium ${j.role === JudgeRole.Main ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-300" : "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300"}`}>
                            {j.role === JudgeRole.Main ? "Главный" : "Боковой"}
                          </span>
                        </div>
                        <div className="flex shrink-0 items-center gap-1">
                          {j.role !== JudgeRole.Main && (
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(j.profile_id, JudgeRole.Main)}>→ Главный</Button>
                          )}
                          {j.role !== JudgeRole.Side && (
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(j.profile_id, JudgeRole.Side)}>→ Боковой</Button>
                          )}
                          <Button size="sm" variant="ghost" className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive" disabled={removingJudge === j.profile_id} onClick={judgeRemove(j.profile_id)}>
                            <X className="h-4 w-4" />
                          </Button>
                        </div>
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </TabsContent>
            <TabsContent value="participants" className="mt-3 space-y-3">
              <div className="space-y-1.5">
                <Label>Поиск участника</Label>
                <Input value={judgeSearchQ} maxLength={50} onChange={(e) => setJudgeSearchQ(e.target.value)} placeholder="Никнейм…" />
              </div>
              <div className="rounded-lg border border-border/40 bg-background/40 p-3">
                {(() => {
                  const nonJudges = (allParticipants.data?.items ?? [])
                    .filter((p) => !(judgesQ.data?.items ?? []).some((j) => j.profile_id === p.id))
                    .filter((p) => !judgeSearchQ.trim() || displayUserName(p).toLowerCase().includes(judgeSearchQ.trim().toLowerCase()));
                  if (!nonJudges.length) return <p className="text-sm text-muted-foreground">{judgeSearchQ.trim() ? "Никого не найдено." : "Все участники уже являются судьями."}</p>;
                  return (
                    <ul className="divide-y divide-border/40">
                      {nonJudges.map((p) => (
                        <li key={p.id} className="flex items-center justify-between gap-2 py-2 text-sm">
                          <span className="truncate">{displayUserName(p)}</span>
                          <div className="flex shrink-0 gap-1">
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(p.id, JudgeRole.Main)}>Главный</Button>
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(p.id, JudgeRole.Side)}>Боковой</Button>
                          </div>
                        </li>
                      ))}
                    </ul>
                  );
                })()}
              </div>
            </TabsContent>
          </Tabs>
        </DialogContent>
      </Dialog>
    </PageShell>
  );
}
