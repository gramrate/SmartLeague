import { useMutation, useQuery } from "@tanstack/react-query";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getClub, getClubMembers, joinClub, leaveClub, setLeader } from "../../api/clubs";
import { queryClient } from "../../shared/queryClient";
import { useAuthStore } from "../../store/authStore";

export function ClubDetailPage() {
  const { id } = useParams();
  const clubId = id!;
  const { userId, role } = useAuthStore();
  const navigate = useNavigate();

  const clubQ = useQuery({ queryKey: ["club", clubId], queryFn: () => getClub(clubId) });
  const membersQ = useQuery({ queryKey: ["club", clubId, "members", { limit: 20, offset: 0 }], queryFn: () => getClubMembers(clubId, { limit: 20, offset: 0 }) });

  const joinM = useMutation({
    mutationFn: () => joinClub(clubId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["club", clubId, "members"] });
      await queryClient.invalidateQueries({ queryKey: ["profile", "me"] });
    }
  });
  const leaveM = useMutation({
    mutationFn: () => leaveClub(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["club", clubId, "members"] });
      await queryClient.invalidateQueries({ queryKey: ["profile", "me"] });
    }
  });

  if (clubQ.isLoading) return <div>Loading...</div>;
  if (clubQ.isError) return <div>Failed to load club</div>;
  if (!clubQ.data) return <div>No data</div>;

  return (
    <div className="space-y-4">
      <div className="rounded bg-white p-6 shadow">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-xl font-semibold">{clubQ.data.name}</h1>
            {clubQ.data.description ? <p className="mt-2 text-sm text-gray-700 whitespace-pre-wrap">{clubQ.data.description}</p> : null}
          </div>
          <div className="flex gap-2">
            <button className="rounded bg-blue-600 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={joinM.isPending} onClick={() => joinM.mutate()}>
              Join
            </button>
            <button className="rounded bg-gray-900 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={leaveM.isPending} onClick={() => leaveM.mutate()}>
              Leave
            </button>
          </div>
        </div>
        <div className="mt-4">
          <Link className="text-sm text-blue-600" to={`/clubs/${clubId}/series`}>
            View series
          </Link>
        </div>
      </div>

      <div className="rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Members</h2>
        {membersQ.isLoading ? <div className="mt-3 text-sm">Loading...</div> : null}
        {membersQ.isError ? <div className="mt-3 text-sm">Failed to load members</div> : null}
        {membersQ.data ? (
          <div className="mt-3 divide-y">
            {membersQ.data.items.map((m) => (
              <div key={m.id} className="flex items-center justify-between py-2">
                <div>
                  <button className="text-left text-sm font-medium text-blue-700 hover:underline" onClick={() => navigate(`/user/${m.id}`)}>
                    {m.name}
                  </button>
                  <div className="text-xs text-gray-600">{m.email}</div>
                </div>
                {role != null && role >= 2 && userId ? (
                  <button
                    className="rounded border px-3 py-1.5 text-xs"
                    onClick={async () => {
                      await setLeader(clubId, m.id);
                      await queryClient.invalidateQueries({ queryKey: ["club", clubId, "members"] });
                    }}
                  >
                    Set leader
                  </button>
                ) : null}
              </div>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  );
}
