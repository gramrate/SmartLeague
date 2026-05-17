import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router-dom";
import { getUserById, getUserGames } from "../../api/users";
import { BackButton } from "../../shared/backButton";

export function UserViewPage() {
  const { id } = useParams();
  const userId = id!;

  const userQ = useQuery({
    queryKey: ["user", userId],
    queryFn: () => getUserById(userId)
  });
  const gamesQ = useQuery({
    queryKey: ["user", userId, "games", { limit: 3, offset: 0 }],
    queryFn: () => getUserGames(userId, { limit: 3, offset: 0 })
  });

  if (userQ.isLoading) return <div>Loading...</div>;

  if (userQ.data) {
    return (
      <div className="space-y-3">
        <BackButton />
        <div className="max-w-xl rounded bg-white p-6 shadow">
          <h1 className="text-xl font-semibold">{userQ.data.name}</h1>
          <div className="mt-1 text-sm text-gray-700">{userQ.data.email}</div>
          <div className="mt-2 text-xs text-gray-600">Nickname: {userQ.data.nickname || "-"}</div>
          <div className="mt-2 text-xs text-gray-600">
            Club: {userQ.data.club_id ? <Link className="text-blue-700 hover:underline" to={`/clubs/${userQ.data.club_id}`}>Open club</Link> : "-"}
          </div>
          <div className="mt-2 text-xs text-gray-600">Club state: {userQ.data.club_state ?? "-"}</div>
          <div className="mt-2 text-xs text-gray-600">Role: {userQ.data.role}</div>
          {userQ.data.description ? <p className="mt-4 whitespace-pre-wrap text-sm">{userQ.data.description}</p> : null}

          <div className="mt-5">
            <div className="mb-2 text-sm font-semibold">Last games</div>
            {gamesQ.data ? (
              <div className="space-y-2">
                {gamesQ.data.items.map((g) => (
                  <Link key={g.id} to={`/game/${g.id}`} className="block rounded border px-3 py-2 text-sm hover:border-gray-400">
                    #{g.number} {g.name || "Untitled game"}
                  </Link>
                ))}
                {gamesQ.data.items.length === 0 ? <div className="text-xs text-gray-600">No games yet</div> : null}
              </div>
            ) : (
              <div className="text-xs text-gray-600">Loading...</div>
            )}
            <div className="mt-3 flex gap-3 text-sm">
              <Link className="text-blue-700 hover:underline" to={`/user/${userId}/games`}>
                View all games
              </Link>
              <Link className="text-blue-700 hover:underline" to={`/user/${userId}/series`}>
                View all series
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return <div>User not found</div>;
}
