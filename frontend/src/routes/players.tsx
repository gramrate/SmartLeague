import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { usersApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { displayUserName } from "@/lib/roles";
import { useDebouncedValue } from "@/lib/useDebouncedValue";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/players")({ component: PlayersPage });

function PlayersPage() {
  const [q, setQ] = useState("");
  const [club, setClub] = useState("");
  const [page, setPage] = useState(1);
  const pageSize = 21;
  const debouncedQ = useDebouncedValue(q, 150);
  const debouncedClub = useDebouncedValue(club, 150);
  const { data, isLoading, error } = useQuery({
    queryKey: ["players", debouncedQ, debouncedClub, page],
    queryFn: () => usersApi.search({
      q: debouncedQ || undefined,
      club: debouncedClub || undefined,
      limit: pageSize,
      offset: (page - 1) * pageSize,
    }),
  });

  return (
    <PageShell>
      <PageHeader eyebrow="Игроки" title="Поиск игрока" description="Поиск по никнейму." />
      <div className="mb-6 grid max-w-3xl gap-3 sm:grid-cols-2">
        <Input
          placeholder="Поиск игроков…"
          value={q}
          onChange={(e) => {
            setQ(e.target.value);
            setPage(1);
          }}
        />
        <Input
          placeholder="Клуб…"
          value={club}
          onChange={(e) => {
            setClub(e.target.value);
            setPage(1);
          }}
        />
      </div>
      {isLoading ? <LoadingBlock /> : error ? <ErrorBlock error={error} /> :
        !data?.items?.length ? <EmptyBlock title="Игроки не найдены" /> : (
        <>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {data.items.map((u) => (
            <Link key={u.id} to="/user/$id" params={{ id: u.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50">
              <p className="font-display text-base font-semibold">{displayUserName(u)}</p>
              {u.description && <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">{u.description}</p>}
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
