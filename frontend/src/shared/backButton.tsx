import { useNavigate } from "react-router-dom";
import { queryClient } from "./queryClient";

export function BackButton() {
  const navigate = useNavigate();

  return (
    <button
      className="inline-flex items-center gap-1 rounded border bg-white px-2 py-1 text-sm hover:bg-gray-50"
      onClick={async () => {
        await queryClient.invalidateQueries();
        navigate(-1);
      }}
    >
      ← Назад
    </button>
  );
}
