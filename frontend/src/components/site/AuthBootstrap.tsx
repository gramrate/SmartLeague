import { useEffect } from "react";
import { useAuthStore } from "@/lib/auth-store";

export function AuthBootstrap() {
  const bootstrap = useAuthStore((s) => s.bootstrap);
  const status = useAuthStore((s) => s.status);
  useEffect(() => {
    if (status === "idle") void bootstrap();
  }, [bootstrap, status]);
  return null;
}
