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
import { toast } from "sonner";
import { Crown, UserX, Ban } from "lucide-react";

export const Route = createFileRoute("/clubs/$id/manage")({ component: ManageClubPage });

function ManageClubPage() {
  const { id } = Route.useParams();
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const navigate = useNavigate();
  const qc = useQueryClient();
  const canManage = canManageClub(me, id);
  const isPresident = me?.club_id === id && me?.club_state === ClubState.President;

  useEffect(() => {
    if (status === "ready" && !canManage) navigate({ to: "/clubs/$id", params: { id } });
  }, [canManage, status, navigate, id]);

  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const members = useQuery({ queryKey: ["club", id, "members"], queryFn: () => clubsApi.members(id) });
  const series = useQuery({ queryKey: ["club", id, "series"], queryFn: () => clubsApi.series(id) });
  const games = useQuery({
    queryKey: ["club", id, "manage-games", series.data?.items?.map((s) => s.id).join(",") ?? ""],
    enabled: !!series.data?.items?.length,
    queryFn: async () => {
      const items = series.data?.items ?? [];
      const responses = await Promise.all(items.map((s) => seriesApi.games(s.id).catch(() => null)));
      return responses.flatMap((res, i) => (res?.items ?? []).map((g) => ({
        ...g,
        _seriesId: items[i].id,
        _seriesName: items[i].name,
      })));
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
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
  };

  const allMembers = members.data?.items ?? [];
  const president = allMembers.find((m) => m.club_state === ClubState.President);
  const others = allMembers.filter((m) => m.id !== president?.id);
  const [presidentTargetId, setPresidentTargetId] = useState<string>("");
  useEffect(() => {
    if (!presidentTargetId && others.length > 0) setPresidentTargetId(others[0].id);
  }, [others, presidentTargetId]);

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
    if (!confirm("Передать президентство этому участнику? Вы станете лидером.")) return;
    try { await clubsApi.setLeader(id, memberId); qc.invalidateQueries({ queryKey: ["club", id, "members"] }); toast.success("Президентство передано"); }
    catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
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
                <Link to="/series/$id" params={{ id: s.id }} className="hover:text-primary">{s.name}</Link>
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
                    <Link to="/series/$id" params={{ id: g._seriesId }} className="hover:underline">
                      {g._seriesName}
                    </Link>
                    {" · #"}
                    {g.number}
                  </p>
                </div>
                <Button size="sm" variant="outline" onClick={() => deleteGame(g.id)}>Удалить</Button>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="mt-8 rounded-2xl border border-destructive/40 bg-card/60 p-6">
        <h2 className="mb-4 font-display text-lg font-semibold text-destructive">Опасные действия</h2>
        <div className="flex flex-wrap items-center gap-3">
          {isPresident && (
            <>
              <Select value={presidentTargetId} onValueChange={setPresidentTargetId}>
                <SelectTrigger className="h-9 w-[260px]">
                  <SelectValue placeholder="Выберите участника" />
                </SelectTrigger>
                <SelectContent>
                  {others.map((m) => (
                    <SelectItem key={m.id} value={m.id}>{displayUserName(m)}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button
                variant="outline"
                disabled={!presidentTargetId}
                onClick={() => presidentTargetId && promoteLeader(presidentTargetId)}
              >
                Передать президентство
              </Button>
            </>
          )}
          <Button
            variant="destructive"
            onClick={async () => {
              if (!confirm("Удалить клуб? Это действие необратимо.")) return;
              try {
                await clubsApi.delete(id);
                toast.success("Клуб удален");
                navigate({ to: "/clubs" });
              } catch (e) {
                toast.error(e instanceof ApiError ? e.message : "Ошибка");
              }
            }}
          >
            Удалить клуб
          </Button>
        </div>
      </section>
    </PageShell>
  );
}
