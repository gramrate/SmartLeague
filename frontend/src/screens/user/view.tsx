import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { getUserById } from "../../api/users";

export function UserViewPage() {
  const { id } = useParams();
  const userId = id!;

  const userQ = useQuery({
    queryKey: ["user", userId],
    queryFn: () => getUserById(userId)
  });

  if (userQ.isLoading) return <div>Loading...</div>;

  if (userQ.data) {
    return (
      <div className="max-w-xl rounded bg-white p-6 shadow">
        <h1 className="text-xl font-semibold">{userQ.data.name}</h1>
        <div className="mt-1 text-sm text-gray-700">{userQ.data.email}</div>
        <div className="mt-2 text-xs text-gray-600">Nickname: {userQ.data.nickname || "-"}</div>
        <div className="mt-2 text-xs text-gray-600">Club state: {userQ.data.club_state ?? "-"}</div>
        <div className="mt-2 text-xs text-gray-600">Role: {userQ.data.role}</div>
        {userQ.data.description ? <p className="mt-4 whitespace-pre-wrap text-sm">{userQ.data.description}</p> : null}
      </div>
    );
  }

  return <div>User not found</div>;
}
