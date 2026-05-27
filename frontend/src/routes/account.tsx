import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useAuthStore } from "@/lib/auth-store";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { RoleBadge } from "@/components/site/RoleBadge";
import { ClubState } from "@/types/api";
import { authApi, clubsApi, ApiError } from "@/lib/api";
import { toast } from "sonner";
import { useQuery } from "@tanstack/react-query";

export const Route = createFileRoute("/account")({ component: AccountPage });

function AccountPage() {
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const setMe = useAuthStore((s) => s.setMe);
  const navigate = useNavigate();

  useEffect(() => { if (status === "ready" && !me) navigate({ to: "/login" }); }, [me, status, navigate]);

  const club = useQuery({
    queryKey: ["club", me?.club_id],
    queryFn: () => clubsApi.get(me!.club_id!),
    enabled: !!me?.club_id,
  });

  const [name, setName] = useState(me?.name ?? "");
  const [nickname, setNickname] = useState(me?.nickname ?? "");
  const [description, setDescription] = useState(me?.description ?? "");
  const [showName, setShowName] = useState(me?.show_name ?? true);
  const [saving, setSaving] = useState(false);
  const nameLimit = 100;
  const nicknameLimit = 100;
  const descriptionLimit = 2000;

  useEffect(() => {
    if (me) {
      setName(me.name ?? ""); setNickname(me.nickname ?? "");
      setDescription(me.description ?? ""); setShowName(me.show_name ?? true);
    }
  }, [me]);

  if (!me) return null;

  const save = async () => {
    setSaving(true);
    try {
      const u = await authApi.updateMe({ name, nickname, description, show_name: showName });
      setMe(u);
      toast.success("Профиль обновлен");
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Не удалось обновить");
    } finally { setSaving(false); }
  };

  return (
    <PageShell>
      <PageHeader eyebrow="Аккаунт" title="Ваш профиль" />
      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-6 lg:col-span-2">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-4 font-display text-lg font-semibold">Данные профиля</h2>
            <div className="grid gap-4 sm:grid-cols-2">
              <Field label="Имя" hint={`${name.length}/${nameLimit} · осталось ${nameLimit - name.length}`}>
                <Input value={name} onChange={(e) => setName(e.target.value)} maxLength={nameLimit} />
              </Field>
              <Field label="Никнейм" hint={`${nickname.length}/${nicknameLimit} · осталось ${nicknameLimit - nickname.length}`}>
                <Input value={nickname} onChange={(e) => setNickname(e.target.value)} maxLength={nicknameLimit} />
              </Field>
            </div>
            <Field label="О себе" className="mt-4" hint={`${description.length}/${descriptionLimit} · осталось ${descriptionLimit - description.length}`}>
              <Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} maxLength={descriptionLimit} />
            </Field>
            <label className="mt-4 flex items-center gap-2 text-sm text-muted-foreground">
              <Checkbox checked={showName} onCheckedChange={(v) => setShowName(!!v)} />
              Показывать мое настоящее имя публично
            </label>
            <div className="mt-6">
              <Button onClick={save} disabled={saving || name.length > nameLimit || nickname.length > nicknameLimit || description.length > descriptionLimit}>{saving ? "Сохранение…" : "Сохранить изменения"}</Button>
            </div>
          </div>
        </div>

        <aside className="space-y-6">
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <p className="text-xs uppercase tracking-widest text-muted-foreground">Почта</p>
            <p className="mt-1 text-sm">{me.email}</p>
          </div>
          <div className="rounded-2xl border border-border/60 bg-card/60 p-6">
            <p className="text-xs uppercase tracking-widest text-muted-foreground">Клуб</p>
            {me.club_id ? (
              <div className="mt-2 space-y-3">
                <Link to="/clubs/$id" params={{ id: me.club_id }} className="block font-semibold text-primary hover:underline">
                  {club.data?.name ?? "Ваш клуб"}
                </Link>
                <RoleBadge state={(me.club_state ?? ClubState.None) as ClubState} />
              </div>
            ) : (
              <div className="mt-2 space-y-2">
                <p className="text-sm text-muted-foreground">Вы пока не состоите в клубе.</p>
                <Button size="sm" asChild><Link to="/clubs">Список клубов</Link></Button>
                <Button size="sm" variant="outline" asChild className="w-full"><Link to="/clubs/create">Создать клуб</Link></Button>
              </div>
            )}
          </div>
        </aside>
      </div>
    </PageShell>
  );
}

function Field({ label, children, className = "", hint }: { label: string; children: React.ReactNode; className?: string; hint?: string }) {
  return (
    <div className={`space-y-1.5 ${className}`}>
      <Label>{label}</Label>
      {children}
      {hint ? <p className="text-xs text-muted-foreground">{hint}</p> : null}
    </div>
  );
}
