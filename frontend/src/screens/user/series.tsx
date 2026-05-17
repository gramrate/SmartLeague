import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router-dom";
import { getUserById, getUserSeries } from "../../api/users";
import { BackButton } from "../../shared/backButton";

export function UserSeriesPage() {
  const { id } = useParams();
  const userId = id!;

  const userQ = useQuery({ queryKey: ["user", userId], queryFn: () => getUserById(userId) });
  const seriesQ = useQuery({ queryKey: ["user", userId, "series", { limit: 100, offset: 0 }], queryFn: () => getUserSeries(userId, { limit: 100, offset: 0 }) });

  return (
    <div className="space-y-4">
      <BackButton />
      <h1 className="text-xl font-semibold">{userQ.data ? `${userQ.data.nickname || userQ.data.name}: series` : "Player series"}</h1>

      {seriesQ.isLoading ? <div>Loading...</div> : null}
      {seriesQ.isError ? <div>Failed to load series</div> : null}

      {seriesQ.data ? (
        <div className="grid gap-3">
          {seriesQ.data.items.map((s) => (
            <Link key={s.id} to={`/series/${s.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
              <div className="font-medium">{s.name}</div>
              <div className="mt-1 text-xs text-gray-600">{new Date(s.start_at).toLocaleString()} - {new Date(s.end_at).toLocaleString()}</div>
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  );
}
