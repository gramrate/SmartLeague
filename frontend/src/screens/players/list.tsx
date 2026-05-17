import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { Link } from "react-router-dom";
import { getUsers } from "../../api/users";

export function PlayersPage() {
  const [search, setSearch] = useState("");

  const q = useQuery({
    queryKey: ["players", { q: search }],
    queryFn: () => getUsers({ limit: 100, offset: 0, q: search.trim() || undefined })
  });

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold">Players</h1>
      <input
        className="w-full rounded border bg-white px-3 py-2 text-sm"
        placeholder="Search by nickname"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />

      {q.isLoading ? <div>Loading...</div> : null}
      {q.isError ? <div>Failed to load players</div> : null}

      {q.data ? (
        <div className="grid gap-3">
          {q.data.items.map((u) => (
            <Link key={u.id} to={`/user/${u.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
              <div className="font-medium">{u.nickname || u.name}</div>
              <div className="mt-1 text-xs text-gray-600">{u.name}</div>
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  );
}
