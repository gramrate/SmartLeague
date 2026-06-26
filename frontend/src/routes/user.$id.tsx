import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { usersApi, clubsApi, authApi, ApiError } from "@/lib/api";
import { translateError } from "@/lib/errors";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { RoleBadge } from "@/components/site/RoleBadge";
import { ClubState } from "@/types/api";
import { displayUserName } from "@/lib/roles";
import { fmtDate, fmtDateRange, fmtRub } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { useAuthStore } from "@/lib/auth-store";
import { toast } from "sonner";
import { useEffect, useState } from "react";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";

export const Route = createFileRoute("/user/$id")({ component: UserPage });

function UserPage() {
  const { id } = Route.useParams();
  const location = useLocation();
  const me = useAuthStore((s) => s.me);
  const setMe = useAuthStore((s) => s.setMe);
  const qc = useQueryClient();
  const isMe = me?.id === id;
  const [editOpen, setEditOpen] = useState(false);
  const [accountOpen, setAccountOpen] = useState(false);
  const [name, setName] = useState("");
  const [nickname, setNickname] = useState("");
  const [description, setDescription] = useState("");
  const [showName, setShowName] = useState(true);
  const [saving, setSaving] = useState(false);
  const [pwdDialogOpen, setPwdDialogOpen] = useState(false);
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [savingPwd, setSavingPwd] = useState(false);
  const [pwdError, setPwdError] = useState<string | null>(null);

  const changePassword = async () => {
    if (oldPassword.length < 8) { setPwdError("Текущий пароль слишком короткий (минимум 8 символов)"); return; }
    if (newPassword.length < 8) { setPwdError("Новый пароль слишком короткий (минимум 8 символов)"); return; }
    setSavingPwd(true);
    setPwdError(null);
    try {
      await authApi.changePassword(oldPassword, newPassword);
      toast.success("Пароль изменён");
      setPwdDialogOpen(false);
      setOldPassword("");
      setNewPassword("");
    } catch (e) {
      setPwdError(translateError(e, "Не удалось сменить пароль"));
    } finally {
      setSavingPwd(false);
    }
  };

  useEffect(() => {
    if (me && isMe) {
      setName(me.name ?? "");
      setNickname(me.nickname ?? "");
      setDescription(me.description ?? "");
      setShowName(me.show_name ?? true);
    }
  }, [me, isMe]);

  const profileDirty = isMe && me && (
    name !== (me.name ?? "") ||
    nickname !== (me.nickname ?? "") ||
    description !== (me.description ?? "") ||
    showName !== (me.show_name ?? true)
  );

  const saveProfile = async () => {
    setSaving(true);
    try {
      const u = await authApi.updateMe({ name, nickname, description, show_name: showName });
      setMe(u);
      qc.invalidateQueries({ queryKey: ["user", id] });
      toast.success("Профиль обновлён");
      setEditOpen(false);
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setSaving(false);
    }
  };
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
          <div className="flex flex-wrap gap-2">
            {isMe && (
              <Button variant={editOpen ? "outline" : "default"} onClick={() => setEditOpen((v) => !v)}>
                {editOpen ? "Отмена" : "Редактировать профиль"}
              </Button>
            )}
            {canBlockFromProfile && (
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
            )}
          </div>
        }
      />

      {editOpen && isMe && (
        <div className="mb-8 rounded-2xl border border-border/60 bg-card/60 p-6">
          <h2 className="mb-4 font-display text-lg font-semibold">Редактирование профиля</h2>
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label>Имя</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} maxLength={50} />
              <p className="text-xs text-muted-foreground">{name.length}/50</p>
            </div>
            <div className="space-y-1.5">
              <Label>Никнейм</Label>
              <Input value={nickname} onChange={(e) => setNickname(e.target.value)} maxLength={50} />
              <p className="text-xs text-muted-foreground">{nickname.length}/50</p>
            </div>
          </div>
          <div className="mt-4 space-y-1.5">
            <Label>О себе</Label>
            <Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} maxLength={500} />
            <p className="text-xs text-muted-foreground">{description.length}/500</p>
          </div>
          <label className="mt-4 flex items-center gap-2 text-sm text-muted-foreground">
            <Checkbox checked={showName} onCheckedChange={(v) => setShowName(!!v)} />
            Показывать моё настоящее имя публично
          </label>
          <div className="mt-6">
            <Button onClick={saveProfile} disabled={saving || !profileDirty}>
              {saving ? "Сохранение…" : "Сохранить изменения"}
            </Button>
          </div>
        </div>
      )}

      <div className="grid gap-6 lg:grid-cols-3">
        <aside className="space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            {isMe && me && (
              <div className="mb-4 border-b border-border/60 pb-4">
                <button
                  className="flex w-full items-center justify-between text-left"
                  onClick={() => setAccountOpen((v) => !v)}
                >
                  <span className="font-display text-lg font-semibold">Аккаунт</span>
                  <span className="text-xs text-muted-foreground">{accountOpen ? "▲" : "▼"}</span>
                </button>
                {accountOpen && (
                  <div className="mt-3 space-y-3">
                    <p className="break-words text-sm"><span className="text-muted-foreground">Почта: </span>{me.email}</p>
                    <Button variant="outline" size="sm" onClick={() => { setOldPassword(""); setNewPassword(""); setPwdError(null); setPwdDialogOpen(true); }}>
                      Сменить пароль
                    </Button>
                  </div>
                )}
              </div>
            )}
            <h2 className="mb-3 font-display text-lg font-semibold">Профиль</h2>
            {user.data.name && user.data.show_name && (
              <p className="break-words text-sm"><span className="text-muted-foreground">Имя: </span>{user.data.name}</p>
            )}
            {user.data.nickname && (
              <p className="break-words text-sm"><span className="text-muted-foreground">Никнейм: </span>{user.data.nickname}</p>
            )}
          </div>
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-3 font-display text-lg font-semibold">Клуб</h2>
            {user.data.club_id ? (
              <div className="space-y-2">
                <Link to="/clubs/$id" params={{ id: user.data.club_id }} className="block break-words font-semibold text-primary hover:underline">
                  {club.data?.name ?? "Открыть клуб"}
                </Link>
                <RoleBadge state={(user.data.club_state ?? ClubState.None) as ClubState} />
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Без клуба</p>
            )}
          </div>
        </aside>

        <AlertDialog open={pwdDialogOpen} onOpenChange={(open) => { if (!open) { setPwdError(null); } setPwdDialogOpen(open); }}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Смена пароля</AlertDialogTitle>
            </AlertDialogHeader>
            <div className="space-y-3 py-2">
              <div className="space-y-1.5">
                <Label>Текущий пароль</Label>
                <Input type="password" value={oldPassword} onChange={(e) => { setOldPassword(e.target.value); setPwdError(null); }} maxLength={100} />
                <p className="text-xs text-muted-foreground">{oldPassword.length}/100</p>
                {oldPassword.length > 0 && oldPassword.length < 8 && <p className="text-xs text-destructive">Минимум 8 символов</p>}
              </div>
              <div className="space-y-1.5">
                <Label>Новый пароль</Label>
                <Input type="password" value={newPassword} onChange={(e) => { setNewPassword(e.target.value); setPwdError(null); }} maxLength={100} />
                <p className="text-xs text-muted-foreground">{newPassword.length}/100</p>
                {newPassword.length > 0 && newPassword.length < 8 && <p className="text-xs text-destructive">Минимум 8 символов</p>}
              </div>
              {pwdError && (
                <div className="rounded-md border border-destructive/40 bg-destructive/10 px-3 py-2 text-sm text-destructive">{pwdError}</div>
              )}
            </div>
            <AlertDialogFooter>
              <AlertDialogCancel disabled={savingPwd}>Отмена</AlertDialogCancel>
              <Button
                onClick={() => void changePassword()}
                disabled={savingPwd || oldPassword.length < 8 || newPassword.length < 8}
              >
                {savingPwd ? "Сохранение…" : "Сменить"}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>

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
                      <p className="break-words font-medium">{g.name || `Игра #${g.number}`}</p>
                      <p className="break-words text-xs text-muted-foreground">{g.series_name}</p>
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
                      <p className="break-words font-medium">{s.name}</p>
                      <p className="text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
                      <div className="mt-1 flex flex-wrap items-center gap-2">
                        {s.is_tournament && (
                          <span className="inline-flex rounded-full bg-purple-100 px-2 py-0.5 text-xs text-purple-800 dark:bg-purple-900/40 dark:text-purple-300">
                            Турнир
                          </span>
                        )}
                        {s.is_rating && !s.is_tournament && (
                          <span className="inline-flex rounded-full bg-sky-100 px-2 py-0.5 text-xs text-sky-800">
                            На рейтинг
                          </span>
                        )}
                        {s.is_club_only && (
                          <span className="inline-flex rounded-full bg-teal-100 px-2 py-0.5 text-xs text-teal-800 dark:bg-teal-900/40 dark:text-teal-300">
                            Для участников клуба
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
