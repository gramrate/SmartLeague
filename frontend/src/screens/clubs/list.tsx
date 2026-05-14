import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { listClubs } from "../../api/clubs";

export function ClubsPage() {
  const q = useQuery({
    queryKey: ["clubs", { limit: 20, offset: 0 }],
    queryFn: () => listClubs({ limit: 20, offset: 0 })
  });

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load clubs</div>;
  if (!q.data) return <div>No data</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">Clubs</h1>
        <Link className="rounded bg-gray-900 px-3 py-2 text-sm text-white" to="/clubs/create">
          Create club
        </Link>
      </div>
      <div className="grid gap-3">
        {q.data.items.map((c) => (
          <Link key={c.id} to={`/clubs/${c.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
            <div className="font-medium">{c.name}</div>
            {c.description ? <div className="mt-1 text-sm text-gray-600 line-clamp-2">{c.description}</div> : null}
          </Link>
        ))}
      </div>
    </div>
  );
}
