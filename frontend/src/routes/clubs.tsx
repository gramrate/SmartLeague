import { createFileRoute, Link, Outlet, useLocation } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuthStore } from "@/lib/auth-store";
import { useState } from "react";
import { Plus } from "lucide-react";
import { useDebouncedValue } from "@/lib/useDebouncedValue";

export const Route = createFileRoute("/clubs")({ component: ClubsPage });

function ClubsPage() {
  const location = useLocation();
  const me = useAuthStore((s) => s.me);
  const [q, setQ] = useState("");
  const qLimit = 100;
  const debouncedQ = useDebouncedValue(q, 150);
  const [page, setPage] = useState(1);
  const pageSize = 15;
  const offset = (page - 1) * pageSize;
  const { data, isLoading, error } = useQuery({
    queryKey: ["clubs", "all", debouncedQ, page],
    queryFn: () => clubsApi.all({ q: debouncedQ || undefined, limit: pageSize, offset }),
  });

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
        <Input
          placeholder="Поиск клубов…"
          value={q}
          maxLength={qLimit}
          onChange={(e) => {
            setQ(e.target.value);
            setPage(1);
          }}
        />
        <p className="mt-1 text-xs text-muted-foreground">{q.length}/{qLimit}</p>
      </div>
      {isLoading ? <LoadingBlock /> : error ? <ErrorBlock error={error} /> :
        !data?.items?.length ? <EmptyBlock title="Клубы не найдены" /> : (
        <>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {data.items.map((c) => (
            <Link key={c.id} to="/clubs/$id" params={{ id: c.id }}
              className="group rounded-xl border border-border/60 bg-card/50 p-5 transition-all hover:border-primary/50 hover:bg-card">
              <h3 className="break-words font-display text-lg font-semibold group-hover:text-primary">{c.name}</h3>
              {c.description && <p className="mt-2 line-clamp-3 text-sm text-muted-foreground">{c.description}</p>}
            </Link>
          ))}
        </div>
        <div className="mt-5 flex items-center justify-between text-sm text-muted-foreground">
          <span>Страница {data.pagination.current_page} из {data.pagination.total_pages}</span>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" disabled={!data.pagination.has_previous || isLoading} onClick={() => setPage((p) => Math.max(1, p - 1))}>Назад</Button>
            <Button variant="outline" size="sm" disabled={!data.pagination.has_next || isLoading} onClick={() => setPage((p) => p + 1)}>Далее</Button>
          </div>
        </div>
        </>
      )}
    </PageShell>
  );
}
