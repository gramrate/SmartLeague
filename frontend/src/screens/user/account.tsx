import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { changePassword, updateMe } from "../../api/users";
import { useAuthStore } from "../../store/authStore";

const updateSchema = z.object({
  nickname: z.string().min(1).max(100).optional(),
  name: z.string().min(2).max(100).optional(),
  show_name: z.boolean().optional(),
  description: z.string().max(2000).optional()
});

const passwordSchema = z
  .object({
    old_password: z.string().min(8).max(100),
    new_password: z.string().min(8).max(100)
  })
  .refine((v) => v.old_password !== v.new_password, { path: ["new_password"], message: "New password must differ" });

type UpdateForm = z.infer<typeof updateSchema>;
type PasswordForm = z.infer<typeof passwordSchema>;

export function AccountPage() {
  const { init, nickname, name, showName, description } = useAuthStore();

  const updateForm = useForm<UpdateForm>({
    resolver: zodResolver(updateSchema),
    mode: "onChange",
    defaultValues: { nickname: nickname ?? "", name: name ?? "", show_name: showName ?? true, description: description ?? "" }
  });

  const passForm = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
    mode: "onChange",
    defaultValues: { old_password: "", new_password: "" }
  });

  const updateM = useMutation({
    mutationFn: (data: UpdateForm) => updateMe({ nickname: data.nickname, name: data.name, show_name: data.show_name, description: data.description ?? null }),
    onSuccess: async () => {
      await init();
    }
  });

  const passM = useMutation({
    mutationFn: (data: PasswordForm) => changePassword(data),
    onSuccess: async () => {
      passForm.reset();
      await init();
    }
  });

  return (
    <div className="space-y-4">
      <div className="max-w-xl rounded bg-white p-6 shadow">
        <h1 className="text-xl font-semibold">Account</h1>
        <form className="mt-4 space-y-3" onSubmit={updateForm.handleSubmit(async (d) => updateM.mutateAsync(d))}>
          <div>
            <label className="text-sm">Name</label>
            <input className="mt-1 w-full rounded border px-3 py-2" {...updateForm.register("name")} />
            {updateForm.formState.errors.name ? <p className="mt-1 text-xs text-red-600">{updateForm.formState.errors.name.message}</p> : null}
          </div>
          <div>
            <label className="text-sm">Nickname</label>
            <input className="mt-1 w-full rounded border px-3 py-2" {...updateForm.register("nickname")} />
          </div>
          <div className="flex items-center gap-2">
            <input type="checkbox" {...updateForm.register("show_name")} />
            <span className="text-sm">Show name</span>
          </div>
          <div>
            <label className="text-sm">Description</label>
            <textarea className="mt-1 w-full rounded border px-3 py-2" rows={3} {...updateForm.register("description")} />
          </div>
          <button className="rounded bg-gray-900 px-4 py-2 text-white disabled:opacity-50" disabled={!updateForm.formState.isValid || updateM.isPending}>
            Save
          </button>
        </form>
      </div>

      <div className="max-w-xl rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Change password</h2>
        <form className="mt-4 space-y-3" onSubmit={passForm.handleSubmit(async (d) => passM.mutateAsync(d))}>
          <div>
            <label className="text-sm">Old password</label>
            <input className="mt-1 w-full rounded border px-3 py-2" type="password" {...passForm.register("old_password")} />
            {passForm.formState.errors.old_password ? <p className="mt-1 text-xs text-red-600">{passForm.formState.errors.old_password.message}</p> : null}
          </div>
          <div>
            <label className="text-sm">New password</label>
            <input className="mt-1 w-full rounded border px-3 py-2" type="password" {...passForm.register("new_password")} />
            {passForm.formState.errors.new_password ? <p className="mt-1 text-xs text-red-600">{passForm.formState.errors.new_password.message}</p> : null}
          </div>
          <button className="rounded bg-blue-600 px-4 py-2 text-white disabled:opacity-50" disabled={!passForm.formState.isValid || passM.isPending}>
            Update password
          </button>
        </form>
      </div>
    </div>
  );
}
