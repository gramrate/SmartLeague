import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi } from "@/lib/api";
import { LoadingBlock, EmptyBlock } from "@/components/site/States";
import { fmtDateRange, fmtRub } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { useState } from "react";

export const Route = createFileRoute("/clubs/$id/series")({ component: ClubSeriesPage });

function ClubSeriesPage() {
  const { id } = Route.useParams();
  const [page, setPage] = useState(1);
  const limit = 20;
  const offset = (page - 1) * limit;
  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const series = useQuery({
    queryKey: ["club", id, "series", page],
    queryFn: () => clubsApi.series(id, limit, offset),
  });

  return (
    <PageShell>
      <PageHeader eyebrow={club.data?.name ?? "Клуб"} title="Все серии" />
      {series.isLoading ? <LoadingBlock /> :
        !series.data?.items?.length ? <EmptyBlock title="Серий нет" /> : (
        <ul className="grid gap-3 sm:grid-cols-2">
          {series.data.items.map((s) => (
            <Link key={s.id} to="/series/$id" params={{ id: s.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50">
              <p className="font-medium">{s.name}</p>
              <p className="mt-1 text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
              <div className="mt-2 flex flex-wrap items-center gap-2">
                <span className="inline-flex rounded-full bg-sky-100 px-2 py-0.5 text-xs text-sky-800">
                  {s.is_rating ? "На рейтинг" : "Без рейтинга"}
                </span>
                {Number(s.price_rub ?? 0) > 0 && (
                  <span className="inline-flex rounded-full bg-amber-100 px-2 py-0.5 text-xs text-amber-800">
                    Платно · {fmtRub(s.price_rub)}
                  </span>
                )}
              </div>
            </Link>
          ))}
        </ul>
      )}
      {series.data?.pagination && (
        <div className="mt-6 flex items-center justify-center gap-3">
          <Button
            variant="outline"
            size="sm"
            disabled={!series.data.pagination.has_previous}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Назад
          </Button>
          <span className="text-sm text-muted-foreground">
            Страница {series.data.pagination.current_page} из {series.data.pagination.total_pages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={!series.data.pagination.has_next}
            onClick={() => setPage((p) => p + 1)}
          >
            Далее
          </Button>
        </div>
      )}
    </PageShell>
  );
}
