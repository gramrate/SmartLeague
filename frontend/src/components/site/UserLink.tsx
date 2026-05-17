import { Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { usersApi } from "@/lib/api";
import { displayUserName } from "@/lib/roles";

export function UserLink({ userId, className }: { userId: string; className?: string }) {
  const { data } = useQuery({
    queryKey: ["user", userId],
    queryFn: () => usersApi.get(userId),
    staleTime: 60_000,
    retry: 0,
  });
  return (
    <Link to="/user/$id" params={{ id: userId }} className={className ?? "text-primary hover:underline"}>
      {displayUserName(data ?? { id: userId })}
    </Link>
  );
}
