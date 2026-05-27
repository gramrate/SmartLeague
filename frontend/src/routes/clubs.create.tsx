import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useAuthStore } from "@/lib/auth-store";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { clubsApi, authApi, ApiError } from "@/lib/api";
import { toast } from "sonner";

export const Route = createFileRoute("/clubs/create")({ component: CreateClubPage });

function CreateClubPage() {
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const setMe = useAuthStore((s) => s.setMe);
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [busy, setBusy] = useState(false);
  const nameLimit = 100;
  const descriptionLimit = 1000;

  useEffect(() => {
    if (status === "ready" && !me) navigate({ to: "/login" });
    else if (me?.club_id) navigate({ to: "/clubs/$id", params: { id: me.club_id } });
  }, [me, status, navigate]);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    try {
      const club = await clubsApi.create({ name: name.trim(), description: description.trim() || undefined });
      const u = await authApi.me().catch(() => null);
      if (u) setMe(u);
      toast.success("Клуб создан");
      navigate({ to: "/clubs/$id", params: { id: club.id } });
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Не удалось выполнить действие");
    } finally { setBusy(false); }
  };

  if (!me || me.club_id) return null;

  return (
    <PageShell>
      <div className="mx-auto max-w-xl">
        <PageHeader eyebrow="Клубы" title="Создать клуб" description="Вы станете президентом нового клуба." />
        <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-border/60 bg-card/60 p-6">
          <div className="space-y-1.5">
            <Label>Название</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} required maxLength={nameLimit} />
            <p className="text-xs text-muted-foreground">{name.length}/{nameLimit}</p>
          </div>
          <div className="space-y-1.5">
            <Label>Описание</Label>
            <Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} maxLength={descriptionLimit} />
            <p className="text-xs text-muted-foreground">{description.length}/{descriptionLimit}</p>
          </div>
          <Button type="submit" disabled={busy || !name.trim() || name.length > nameLimit || description.length > descriptionLimit}>{busy ? "Создаем…" : "Создать клуб"}</Button>
        </form>
      </div>
    </PageShell>
  );
}
