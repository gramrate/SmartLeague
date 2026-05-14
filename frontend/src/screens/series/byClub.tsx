import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router-dom";
import { getClubSeries } from "../../api/series";

export function ClubSeriesPage() {
  const { id } = useParams();
  const clubId = id!;

  const q = useQuery({
    queryKey: ["club", clubId, "series", { limit: 20, offset: 0 }],
    queryFn: () => getClubSeries(clubId, { limit: 20, offset: 0 })
  });

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load series</div>;
  if (!q.data) return <div>No data</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">Club series</h1>
        <Link className="rounded bg-gray-900 px-3 py-2 text-sm text-white" to="/series/create">
          Create series
        </Link>
      </div>

      <div className="grid gap-3">
        {q.data.items.map((s) => (
          <Link key={s.id} to={`/series/${s.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
            <div className="font-medium">{s.name}</div>
            <div className="mt-1 text-xs text-gray-600">
              {new Date(s.start_at).toLocaleString()} → {new Date(s.end_at).toLocaleString()}
            </div>
            {s.description ? <div className="mt-2 text-sm text-gray-700 line-clamp-2">{s.description}</div> : null}
          </Link>
        ))}
      </div>
    </div>
  );
}
