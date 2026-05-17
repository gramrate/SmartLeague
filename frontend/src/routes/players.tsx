import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { usersApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { displayUserName } from "@/lib/roles";

export const Route = createFileRoute("/players")({ component: PlayersPage });

function PlayersPage() {
  const [q, setQ] = useState("");
  const { data, isLoading, error } = useQuery({
    queryKey: ["players", q],
    queryFn: () => usersApi.search({ q: q || undefined, limit: 50 }),
  });

  return (
    <PageShell>
      <PageHeader eyebrow="Игроки" title="Поиск игрока" description="Поиск по никнейму, имени или email." />
      <div className="mb-6 max-w-md">
        <Input placeholder="Поиск игроков…" value={q} onChange={(e) => setQ(e.target.value)} />
      </div>
      {isLoading ? <LoadingBlock /> : error ? <ErrorBlock error={error} /> :
        !data?.items?.length ? <EmptyBlock title="Игроки не найдены" /> : (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {data.items.map((u) => (
            <Link key={u.id} to="/user/$id" params={{ id: u.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50">
              <p className="font-display text-base font-semibold">{displayUserName(u)}</p>
              {u.description && <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">{u.description}</p>}
            </Link>
          ))}
        </div>
      )}
    </PageShell>
  );
}
