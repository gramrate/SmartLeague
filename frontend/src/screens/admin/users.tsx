import { useQuery } from "@tanstack/react-query";
import { getUsers } from "../../api/users";

export function UsersAdminPage() {
  const q = useQuery({ queryKey: ["admin", "users", { limit: 20, offset: 0 }], queryFn: () => getUsers({ limit: 20, offset: 0 }) });

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load users</div>;
  if (!q.data) return <div>No data</div>;

  return (
    <div className="rounded bg-white p-6 shadow">
      <h1 className="text-xl font-semibold">Users</h1>
      <div className="mt-4 overflow-auto">
        <table className="min-w-full text-left text-sm">
          <thead>
            <tr className="border-b text-xs text-gray-600">
              <th className="py-2 pr-4">ID</th>
              <th className="py-2 pr-4">Email</th>
              <th className="py-2 pr-4">Name</th>
              <th className="py-2 pr-4">Role</th>
            </tr>
          </thead>
          <tbody>
            {q.data.items.map((u) => (
              <tr key={u.id} className="border-b">
                <td className="py-2 pr-4">{u.id}</td>
                <td className="py-2 pr-4">{u.email}</td>
                <td className="py-2 pr-4">
                  {u.name}
                </td>
                <td className="py-2 pr-4">{u.role}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
