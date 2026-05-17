import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuthStore } from "@/lib/auth-store";
import { useState, useMemo } from "react";
import { Plus } from "lucide-react";
import { includesCaseInsensitive } from "@/lib/search";

export const Route = createFileRoute("/clubs")({ component: ClubsPage });

function ClubsPage() {
  const location = useLocation();
  const me = useAuthStore((s) => s.me);
  const [q, setQ] = useState("");
  const { data, isLoading, error } = useQuery({ queryKey: ["clubs", "all"], queryFn: () => clubsApi.all(200, 0) });

  const items = useMemo(() => {
    const xs = data?.items ?? [];
    if (!q.trim()) return xs;
    return xs.filter((c) => includesCaseInsensitive(c.name, q) || includesCaseInsensitive(c.description, q));
  }, [data, q]);

  if (location.pathname !== "/clubs") {
    return <Outlet />;
  }

  return (
    <PageShell>
      <PageHeader
        eyebrow="Клубы" title="Все клубы"
        description="Открывайте клубы, смотрите участников, серии и игры."
        actions={
          me && !me.club_id ? (
            <Button asChild><Link to="/clubs/create"><Plus className="mr-1 h-4 w-4" />Создать клуб</Link></Button>
          ) : me?.club_id ? (
            <Button variant="outline" asChild><Link to="/clubs/$id" params={{ id: me.club_id }}>Мой клуб</Link></Button>
          ) : null
        }
      />
      <div className="mb-6 max-w-md">
        <Input placeholder="Поиск клубов…" value={q} onChange={(e) => setQ(e.target.value)} />
      </div>
      {isLoading ? <LoadingBlock /> : error ? <ErrorBlock error={error} /> :
        items.length === 0 ? <EmptyBlock title="Клубы не найдены" /> : (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {items.map((c) => (
            <Link key={c.id} to="/clubs/$id" params={{ id: c.id }}
              className="group rounded-xl border border-border/60 bg-card/50 p-5 transition-all hover:border-primary/50 hover:bg-card">
              <h3 className="font-display text-lg font-semibold group-hover:text-primary">{c.name}</h3>
              {c.description && <p className="mt-2 line-clamp-3 text-sm text-muted-foreground">{c.description}</p>}
            </Link>
          ))}
        </div>
      )}
    </PageShell>
  );
}
