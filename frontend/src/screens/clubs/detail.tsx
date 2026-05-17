import { useMutation, useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getClub, getClubMembers, joinClub, leaveClub } from "../../api/clubs";
import { listGames } from "../../api/games";
import { getClubSeries } from "../../api/series";
import { queryClient } from "../../shared/queryClient";
import { useAuthStore } from "../../store/authStore";
import { BackButton } from "../../shared/backButton";
import { ClubState, GameStatus } from "../../types/enums";

function clubStateLabel(state: ClubState) {
  switch (state) {
    case ClubState.Member:
      return "Member";
    case ClubState.Leader:
      return "Leader";
    case ClubState.President:
      return "President";
    case ClubState.Resident:
      return "Resident";
    default:
      return "None";
  }
}

export function ClubDetailPage() {
  const { id } = useParams();
  const clubId = id!;
  const { userId, clubId: myClubId, clubState, init } = useAuthStore();
  const navigate = useNavigate();
  const canManageClub = myClubId === clubId && (clubState === ClubState.Leader || clubState === ClubState.President);
  const [joinError, setJoinError] = useState<string | null>(null);
  const isInAnotherClub = !!userId && !!myClubId && myClubId !== clubId;

  const clubQ = useQuery({ queryKey: ["club", clubId], queryFn: () => getClub(clubId) });
  const membersQ = useQuery({ queryKey: ["club", clubId, "members", { limit: 50, offset: 0 }], queryFn: () => getClubMembers(clubId, { limit: 50, offset: 0 }) });
  const seriesQ = useQuery({ queryKey: ["club", clubId, "series", { limit: 20, offset: 0 }], queryFn: () => getClubSeries(clubId, { limit: 20, offset: 0 }) });

  const currentGamesQ = useQuery({
    queryKey: ["club", clubId, "games", "current"],
    queryFn: async () => {
      const seriesResp = await getClubSeries(clubId, { limit: 20, offset: 0 });
      const perSeries = await Promise.all(
        seriesResp.items.map(async (series) => {
          const gamesResp = await listGames(series.id, { limit: 20, offset: 0 });
          return gamesResp.items.map((g) => ({ ...g, series_name: series.name }));
        })
      );

      return perSeries.flat().filter((g) => g.status !== GameStatus.Finished).sort((a, b) => a.status - b.status || b.number - a.number).slice(0, 6);
    }
  });

  const joinM = useMutation({
    mutationFn: () => joinClub(clubId),
    onSuccess: async () => {
      setJoinError(null);
      await queryClient.invalidateQueries({ queryKey: ["club", clubId, "members"] });
      await queryClient.invalidateQueries({ queryKey: ["profile", "me"] });
      await init();
    },
    onError: (err: any) => {
      const message = String(err?.message || "");
      if (message.includes("already in another club")) {
        setJoinError("Вы уже состоите в другом клубе. Сначала выйдите из него.");
        return;
      }
      if (message.includes("blocked")) {
        setJoinError("Вы заблокированы в этом клубе.");
        return;
      }
      setJoinError("Не удалось вступить в клуб");
    }
  });
  const leaveM = useMutation({
    mutationFn: () => leaveClub(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["club", clubId, "members"] });
      await queryClient.invalidateQueries({ queryKey: ["profile", "me"] });
      await init();
    }
  });

  if (clubQ.isLoading) return <div>Loading...</div>;
  if (clubQ.isError) return <div>Failed to load club</div>;
  if (!clubQ.data) return <div>No data</div>;

  return (
    <div className="space-y-4">
      <BackButton />
      <div className="rounded bg-white p-6 shadow">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-xl font-semibold">{clubQ.data.name}</h1>
            {clubQ.data.description ? <p className="mt-2 text-sm text-gray-700 whitespace-pre-wrap">{clubQ.data.description}</p> : null}
          </div>
          <div className="flex gap-2">
            <button className="rounded bg-blue-600 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={joinM.isPending || isInAnotherClub} onClick={() => joinM.mutate()}>
              Join
            </button>
            <button className="rounded bg-gray-900 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={leaveM.isPending} onClick={() => leaveM.mutate()}>
              Leave
            </button>
          </div>
        </div>
        {canManageClub ? (
          <Link className="mt-3 inline-block rounded bg-gray-900 px-3 py-2 text-sm text-white" to={`/clubs/${clubId}/manage`}>
            Management panel
          </Link>
        ) : null}
        {isInAnotherClub ? <p className="mt-2 text-sm text-red-600">Вы уже состоите в другом клубе. Сначала выйдите из текущего клуба.</p> : null}
        {joinError ? <p className="mt-2 text-sm text-red-600">{joinError}</p> : null}
      </div>

      <div className="rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Members</h2>
        {membersQ.data ? (
          <div className="mt-3 divide-y">
            {membersQ.data.items.map((m) => (
              <div key={m.id} className="flex items-center justify-between py-2">
                <div>
                  <button className="text-left text-sm font-medium text-blue-700 hover:underline" onClick={() => navigate(`/user/${m.id}`)}>
                    {m.name}
                  </button>
                  <div className="text-xs text-gray-600">{m.email}</div>
                  <div className="text-xs text-gray-600">Role in club: {clubStateLabel(m.club_state)}</div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="mt-3 text-sm text-gray-600">Loading...</div>
        )}
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div className="rounded bg-white p-6 shadow">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-lg font-semibold">Series</h2>
            {canManageClub ? <Link className="rounded bg-gray-900 px-3 py-2 text-sm text-white" to="/series/create">Create series</Link> : null}
          </div>
          {seriesQ.data ? (
            <div className="space-y-2">
              {seriesQ.data.items.map((s) => (
                <div key={s.id} className="rounded border px-3 py-2">
                  <Link to={`/series/${s.id}`} className="font-medium hover:underline">{s.name}</Link>
                  <div className="mt-1 text-xs text-gray-600">{new Date(s.start_at).toLocaleString()} - {new Date(s.end_at).toLocaleString()}</div>
                </div>
              ))}
            </div>
          ) : null}
        </div>

        <div className="rounded bg-white p-6 shadow">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold">Current games</h2>
            <Link className="text-sm text-blue-600" to={`/clubs/${clubId}/games`}>View all games</Link>
          </div>
          {currentGamesQ.data ? (
            <div className="mt-3 grid gap-2">
              {currentGamesQ.data.map((g) => (
                <div key={g.id} className="rounded border px-3 py-2">
                  <Link to={`/game/${g.id}`} className="font-medium hover:underline">#{g.number} {g.name || "Untitled game"}</Link>
                  <div className="text-xs text-gray-600">Series: {g.series_name}</div>
                </div>
              ))}
              {currentGamesQ.data.length === 0 ? <div className="text-sm text-gray-600">No current games</div> : null}
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}
