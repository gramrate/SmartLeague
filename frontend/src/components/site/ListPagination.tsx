import { Button } from "@/components/ui/button";
import type { PaginationInfo } from "@/types/api";
import { cn } from "@/lib/utils";

type ListPaginationProps = {
  pagination: PaginationInfo;
  isLoading?: boolean;
  className?: string;
  onPrevious: () => void;
  onNext: () => void;
};

export function ListPagination({
  pagination,
  isLoading = false,
  className,
  onPrevious,
  onNext,
}: ListPaginationProps) {
  return (
    <div
      className={cn(
        "mt-6 flex flex-wrap items-center justify-center gap-3 text-sm text-muted-foreground",
        className,
      )}
    >
      <Button
        variant="outline"
        size="sm"
        disabled={!pagination.has_previous || isLoading}
        onClick={onPrevious}
      >
        Назад
      </Button>
      <span>
        Страница {pagination.current_page} из {pagination.total_pages}
      </span>
      <Button
        variant="outline"
        size="sm"
        disabled={!pagination.has_next || isLoading}
        onClick={onNext}
      >
        Далее
      </Button>
    </div>
  );
}
