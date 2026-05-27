import { createFileRoute, Link } from "@tanstack/react-router";
import { PageHeader, PageShell } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { usersApi } from "@/lib/api";
import { EmptyBlock, LoadingBlock } from "@/components/site/States";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { fmtDate } from "@/lib/format";

export const Route = createFileRoute("/user/$id/games")({ component: UserGamesPage });

function UserGamesPage() {
  const { id } = Route.useParams();
  const [page, setPage] = useState(1);
  const limit = 15;
  const offset = (page - 1) * limit;

  const user = useQuery({ queryKey: ["user", id], queryFn: () => usersApi.get(id) });
  const games = useQuery({
    queryKey: ["user", id, "games", page],
    queryFn: () => usersApi.games(id, limit, offset),
  });

  return (
    <PageShell>
      <PageHeader
        eyebrow={user.data?.nickname ?? "Игрок"}
        title="Все игры игрока"
        actions={<Button variant="outline" asChild><Link to="/user/$id" params={{ id }}>К профилю</Link></Button>}
      />

      {games.isLoading ? <LoadingBlock /> :
        !games.data?.items?.length ? <EmptyBlock title="Игр нет" /> : (
          <ul className="space-y-2">
            {games.data.items.map((g) => (
              <Link
                key={g.id}
                to="/game/$id"
                params={{ id: g.id }}
                className="flex items-center justify-between gap-3 rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50"
              >
                <div className="min-w-0">
                  <p className="truncate font-medium">{g.name || `Игра #${g.number}`}</p>
                  <p className="text-xs text-muted-foreground">{g.series_name}</p>
                </div>
                <span className="shrink-0 text-xs text-muted-foreground">{fmtDate(g.created_at)}</span>
              </Link>
            ))}
          </ul>
      )}

      {games.data?.pagination && (
        <div className="mt-6 flex items-center justify-center gap-3">
          <Button
            variant="outline"
            size="sm"
            disabled={!games.data.pagination.has_previous}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Назад
          </Button>
          <span className="text-sm text-muted-foreground">
            Страница {games.data.pagination.current_page} из {games.data.pagination.total_pages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={!games.data.pagination.has_next}
            onClick={() => setPage((p) => p + 1)}
          >
            Далее
          </Button>
        </div>
      )}
    </PageShell>
  );
}
