import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { clubsApi, ApiError, seriesApi, gamesApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock } from "@/components/site/States";
import { useAuthStore } from "@/lib/auth-store";
import { canManageClub, displayUserName, CLUB_STATE_LABEL } from "@/lib/roles";
import { ClubState } from "@/types/api";
import { RoleBadge } from "@/components/site/RoleBadge";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { toast } from "sonner";
import { Crown, UserX, Ban } from "lucide-react";
import { fmtDateRange } from "@/lib/format";

export const Route = createFileRoute("/clubs/$id/manage")({ component: ManageClubPage });

function ManageClubPage() {
  const { id } = Route.useParams();
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const refreshMe = useAuthStore((s) => s.refreshMe);
  const navigate = useNavigate();
  const qc = useQueryClient();
  const canManage = canManageClub(me, id);
  const isPresident = me?.club_id === id && me?.club_state === ClubState.President;

  useEffect(() => {
    if (status === "ready" && !canManage) navigate({ to: "/clubs/$id", params: { id } });
  }, [canManage, status, navigate, id]);

  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const members = useQuery({ queryKey: ["club", id, "members"], queryFn: () => clubsApi.members(id) });
  const series = useQuery({
    queryKey: ["club", id, "series", "all"],
    queryFn: async () => {
      const limit = 200;
      let offset = 0;
      const all: any[] = [];
      for (;;) {
        const page = await clubsApi.series(id, limit, offset);
        all.push(...(page.items ?? []));
        if (!page.pagination?.has_next) break;
        offset += limit;
      }
      all.sort((a, b) => new Date(a.start_at).getTime() - new Date(b.start_at).getTime());
      return { items: all };
    },
  });
  const games = useQuery({
    queryKey: ["club", id, "manage-games", series.data?.items?.map((s) => s.id).join(",") ?? ""],
    enabled: !!series.data?.items?.length,
    queryFn: async () => {
      const items = series.data?.items ?? [];
      const seriesGames = await Promise.all(items.map(async (s) => {
        const limit = 200;
        let offset = 0;
        const allGames: any[] = [];
        for (;;) {
          const page = await seriesApi.games(s.id, limit, offset).catch(() => null);
          if (!page) break;
          allGames.push(...(page.items ?? []));
          if (!page.pagination?.has_next) break;
          offset += limit;
        }
        return allGames.map((g) => ({
          ...g,
          _seriesId: s.id,
          _seriesName: s.name,
        }));
      }));
      return seriesGames.flat();
    },
  });

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  useEffect(() => {
    if (club.data) { setName(club.data.name); setDescription(club.data.description ?? ""); }
  }, [club.data]);

  if (!canManage) return null;
  if (club.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (club.error) return <PageShell><ErrorBlock error={club.error} /></PageShell>;

  const saveClub = async () => {
    try {
      await clubsApi.update(id, { name, description });
      qc.invalidateQueries({ queryKey: ["club", id] });
      toast.success("Клуб обновлен");
      navigate({ to: "/clubs/$id", params: { id } });
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };

  const allMembers = members.data?.items ?? [];
  const president = allMembers.find((m) => m.club_state === ClubState.President);
  const others = allMembers.filter((m) => m.id !== president?.id);
  const [presidentTargetId, setPresidentTargetId] = useState<string>("");
  const [presidentTargetQuery, setPresidentTargetQuery] = useState("");
  const [presidentDropdownOpen, setPresidentDropdownOpen] = useState(false);
  const [transferOpen, setTransferOpen] = useState(false);
  const [transferring, setTransferring] = useState(false);
  const [deletingClub, setDeletingClub] = useState(false);
  useEffect(() => {
    if (!transferOpen) return;
    setPresidentTargetId("");
    setPresidentTargetQuery("");
  }, [transferOpen]);

  const setRole = async (memberId: string, state: ClubState) => {
    try {
      await clubsApi.setRole(id, memberId, state);
      qc.invalidateQueries({ queryKey: ["club", id, "members"] });
      toast.success("Роль обновлена");
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };
  const kick = async (memberId: string) => {
    if (!confirm("Исключить этого участника?")) return;
    try { await clubsApi.kick(id, memberId); qc.invalidateQueries({ queryKey: ["club", id, "members"] }); toast.success("Участник исключен"); }
    catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };
  const block = async (memberId: string) => {
    if (!confirm("Заблокировать этого участника?")) return;
    try { await clubsApi.block(id, memberId); qc.invalidateQueries({ queryKey: ["club", id, "members"] }); toast.success("Участник заблокирован"); }
    catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };
  const promoteLeader = async (memberId: string) => {
    if (!confirm("Передать президентство выбранному участнику? Вы станете лидером.")) return;
    setTransferring(true);
    try {
      await clubsApi.setLeader(id, memberId);
      await refreshMe();
      qc.invalidateQueries({ queryKey: ["club", id] });
      qc.invalidateQueries({ queryKey: ["club", id, "members"] });
      qc.invalidateQueries({ queryKey: ["club", id, "series"] });
      qc.invalidateQueries({ queryKey: ["club", id, "manage-games"] });
      toast.success("Президентство передано");
      setTransferOpen(false);
    }
    catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
    finally { setTransferring(false); }
  };
  const deleteSeries = async (sid: string) => {
    if (!confirm("Удалить эту серию?")) return;
    try { await (await import("@/lib/api")).seriesApi.delete(sid); qc.invalidateQueries({ queryKey: ["club", id, "series"] }); toast.success("Удалено"); }
    catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };
  const deleteGame = async (gid: string) => {
    if (!confirm("Удалить эту игру?")) return;
    try {
      await gamesApi.delete(gid);
      qc.invalidateQueries({ queryKey: ["club", id, "manage-games"] });
      toast.success("Удалено");
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };

  return (
    <PageShell>
      <PageHeader eyebrow="Управление" title={club.data!.name} description="Редактирование клуба, участников и серий." />

      {/* Edit club */}
      <section className="mb-8 rounded-2xl border border-border/60 bg-card/60 p-6">
        <h2 className="mb-4 font-display text-lg font-semibold">Данные клуба</h2>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-1.5"><Label>Название</Label><Input value={name} onChange={(e) => setName(e.target.value)} /></div>
        </div>
        <div className="mt-4 space-y-1.5"><Label>Описание</Label><Textarea rows={3} value={description} onChange={(e) => setDescription(e.target.value)} /></div>
        <Button className="mt-4" onClick={saveClub}>Сохранить</Button>
      </section>

      {/* President */}
      {president && (
        <section className="mb-8 rounded-2xl border border-primary/40 bg-gradient-to-br from-primary/15 to-card/60 p-6 shadow-[var(--shadow-glow)]">
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <Crown className="h-6 w-6 text-primary" />
              <div>
                <p className="text-xs uppercase tracking-widest text-primary">Президент</p>
                <Link to="/user/$id" params={{ id: president.id }} className="font-display text-xl font-bold hover:underline">
                  {displayUserName(president)}
                </Link>
              </div>
            </div>
            <RoleBadge state={ClubState.President} />
          </div>
        </section>
      )}

      {/* Members */}
      <section className="mb-8 rounded-2xl border border-border/60 bg-card/60 p-6">
        <h2 className="mb-4 font-display text-lg font-semibold">Участники</h2>
        {others.length === 0 ? <p className="text-sm text-muted-foreground">Других участников нет.</p> : (
          <ul className="divide-y divide-border/40">
            {others.map((m) => {
              const state = (m.club_state ?? ClubState.Member) as ClubState;
              return (
                <li key={m.id} className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
                  <div className="flex items-center gap-3">
                    <Link to="/user/$id" params={{ id: m.id }} className="font-medium hover:text-primary">{displayUserName(m)}</Link>
                    <RoleBadge state={state} />
                  </div>
                  <div className="flex flex-wrap items-center gap-2">
                    <Select value={String(state)} onValueChange={(v) => setRole(m.id, Number(v) as ClubState)}>
                      <SelectTrigger className="h-8 w-[140px]"><SelectValue /></SelectTrigger>
                      <SelectContent>
                        <SelectItem value={String(ClubState.Member)}>{CLUB_STATE_LABEL[ClubState.Member]}</SelectItem>
                        <SelectItem value={String(ClubState.Resident)}>{CLUB_STATE_LABEL[ClubState.Resident]}</SelectItem>
                        <SelectItem value={String(ClubState.Leader)}>{CLUB_STATE_LABEL[ClubState.Leader]}</SelectItem>
                      </SelectContent>
                    </Select>
                    <Button size="sm" variant="outline" onClick={() => kick(m.id)}><UserX className="h-4 w-4" /></Button>
                    <Button size="sm" variant="outline" onClick={() => block(m.id)}><Ban className="h-4 w-4" /></Button>
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </section>

      {/* Series management */}
      <section className="rounded-2xl border border-border/60 bg-card/60 p-6">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="font-display text-lg font-semibold">Серии</h2>
          <Button size="sm" asChild><Link to="/series/create">Создать серию</Link></Button>
        </div>
        {!series.data?.items?.length ? <p className="text-sm text-muted-foreground">Серий нет.</p> : (
          <ul className="divide-y divide-border/40">
            {series.data.items.map((s) => (
              <li key={s.id} className="flex items-center justify-between py-3">
                <div>
                  <Link to="/series/$id" params={{ id: s.id }} className="hover:text-primary">{s.name}</Link>
                  <p className="text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
                  <p className="text-xs text-muted-foreground">{s.is_rating ? "На рейтинг" : "Без рейтинга"}</p>
                </div>
                <Button size="sm" variant="outline" onClick={() => deleteSeries(s.id)}>Удалить</Button>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="mt-8 rounded-2xl border border-border/60 bg-card/60 p-6">
        <h2 className="mb-4 font-display text-lg font-semibold">Игры</h2>
        {games.isLoading ? <LoadingBlock /> : !games.data?.length ? (
          <p className="text-sm text-muted-foreground">Игр нет.</p>
        ) : (
          <ul className="divide-y divide-border/40">
            {games.data.map((g) => (
              <li key={g.id} className="flex items-center justify-between py-3">
                <div>
                  <Link to="/game/$id" params={{ id: g.id }} className="hover:text-primary">
                    {g.name || `Игра #${g.number}`}
                  </Link>
                  <p className="text-xs text-muted-foreground">
                    Серия:
                    {" "}
                    <Link to="/series/$id" params={{ id: g._seriesId }} className="hover:underline">
                      {g._seriesName}
                    </Link>
                    {" · Игра #"}
                    {g.number}
                  </p>
                </div>
                <Button size="sm" variant="outline" onClick={() => deleteGame(g.id)}>Удалить</Button>
              </li>
            ))}
          </ul>
        )}
      </section>

      {isPresident && (
        <section className="mt-8 rounded-2xl border border-destructive/40 bg-card/60 p-6">
          <h2 className="mb-4 font-display text-lg font-semibold text-destructive">Опасные действия</h2>
          <div className="flex flex-wrap items-center gap-3">
            <Button variant="outline" onClick={() => setTransferOpen(true)} disabled={!others.length}>
              Передать президентство
            </Button>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="destructive">Удалить клуб</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Удалить клуб?</AlertDialogTitle>
                  <AlertDialogDescription>
                    Это действие необратимо.
                    {" "}
                    При удалении клуба президентство передастся случайному лидеру, если лидеров нет — случайному резиденту, если резидентов нет — случайному участнику.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel>Отмена</AlertDialogCancel>
                  <AlertDialogAction
                    disabled={deletingClub}
                    onClick={async (e) => {
                      e.preventDefault();
                      setDeletingClub(true);
                      try {
                        await clubsApi.delete(id);
                        toast.success("Клуб удален");
                        navigate({ to: "/clubs" });
                      } catch (err) {
                        toast.error(err instanceof ApiError ? err.message : "Ошибка");
                      } finally {
                        setDeletingClub(false);
                      }
                    }}
                  >
                    {deletingClub ? "Удаление..." : "Удалить"}
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </div>
        </section>
      )}

      <Dialog open={transferOpen} onOpenChange={setTransferOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Передать президентство</DialogTitle>
            <DialogDescription>
              Выберите участника клуба и нажмите кнопку передачи.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-2">
            <Label>Новый президент</Label>
            <div className="relative">
              <Input
                value={presidentTargetQuery}
                onChange={(e) => {
                  const value = e.target.value;
                  setPresidentTargetQuery(value);
                  setPresidentDropdownOpen(true);
                  const exact = others.find((m) => displayUserName(m).toLowerCase() === value.trim().toLowerCase());
                  setPresidentTargetId(exact?.id ?? "");
                }}
                onFocus={() => setPresidentDropdownOpen(true)}
                onBlur={() => setTimeout(() => setPresidentDropdownOpen(false), 120)}
                placeholder="Введите или выберите участника"
                className="h-9"
              />
              {presidentDropdownOpen && (
                <div className="absolute z-20 mt-1 max-h-56 w-full overflow-auto rounded-md border border-border bg-popover p-1 shadow-md">
                  {others
                    .filter((m) => displayUserName(m).toLowerCase().includes(presidentTargetQuery.trim().toLowerCase()))
                    .map((m) => (
                      <button
                        key={m.id}
                        type="button"
                        className="block w-full rounded-sm px-2 py-1.5 text-left text-sm hover:bg-accent"
                        onMouseDown={(e) => e.preventDefault()}
                        onClick={() => {
                          setPresidentTargetId(m.id);
                          setPresidentTargetQuery(displayUserName(m));
                          setPresidentDropdownOpen(false);
                        }}
                      >
                        {displayUserName(m)}
                      </button>
                    ))}
                </div>
              )}
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setTransferOpen(false)}>Отмена</Button>
            <Button disabled={!presidentTargetId || transferring} onClick={() => presidentTargetId && promoteLeader(presidentTargetId)}>
              {transferring ? "Передача..." : "Передать"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </PageShell>
  );
}
