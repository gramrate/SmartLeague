import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi, seriesApi } from "@/lib/api";
import { LoadingBlock, EmptyBlock } from "@/components/site/States";
import { fmtDateRange } from "@/lib/format";

export const Route = createFileRoute("/clubs/$id/series")({ component: ClubSeriesPage });

function ClubSeriesPage() {
  const { id } = Route.useParams();
  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const series = useQuery({ queryKey: ["club", id, "series"], queryFn: () => clubsApi.series(id) });

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
            </Link>
          ))}
        </ul>
      )}
    </PageShell>
  );
}
