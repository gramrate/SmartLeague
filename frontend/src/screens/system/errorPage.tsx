import { isRouteErrorResponse, Link, useRouteError } from "react-router-dom";

export function ErrorPage() {
  const error = useRouteError();

  let title = "Unexpected error";
  let message = "Something went wrong";

  if (isRouteErrorResponse(error)) {
    title = `${error.status} ${error.statusText}`;
    message = typeof error.data === "string" ? error.data : message;
  } else if (error instanceof Error) {
    message = error.message;
  }

  return (
    <div className="mx-auto max-w-xl rounded bg-white p-8 shadow">
      <h1 className="text-2xl font-semibold">{title}</h1>
      <p className="mt-3 text-sm text-gray-700">{message}</p>
      <div className="mt-6 flex gap-3">
        <Link to="/" className="rounded bg-gray-900 px-4 py-2 text-sm text-white">
          Go home
        </Link>
        <button className="rounded border px-4 py-2 text-sm" onClick={() => window.location.reload()}>
          Reload
        </button>
      </div>
    </div>
  );
}

