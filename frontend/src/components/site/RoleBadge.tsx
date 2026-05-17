import { cn } from "@/lib/utils";
import { CLUB_STATE_LABEL, CLUB_STATE_TONE } from "@/lib/roles";
import type { ClubState } from "@/types/api";

export function RoleBadge({ state, className }: { state: ClubState; className?: string }) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-semibold",
        CLUB_STATE_TONE[state],
        className,
      )}
    >
      {CLUB_STATE_LABEL[state]}
    </span>
  );
}
