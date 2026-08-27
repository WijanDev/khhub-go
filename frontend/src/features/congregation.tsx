import { useEffect, useState, type FormEvent } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Card, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select";
import { api, type Congregation } from "@/lib/api";
import { weekdayLabels } from "@/lib/labels";

export function CongregationPage() {
  const qc = useQueryClient();
  const q = useQuery({
    queryKey: ["congregation"],
    queryFn: () => api<Congregation>("/congregation"),
  });
  const [form, setForm] = useState<Congregation>({
    name: "",
    number: "",
    midweekDay: 4,
    weekendDay: 0,
  });
  const [pw, setPw] = useState({ currentPassword: "", newPassword: "" });
  const [msg, setMsg] = useState("");
  const [resetOpen, setResetOpen] = useState(false);

  useEffect(() => {
    if (q.data) setForm(q.data);
  }, [q.data]);

  const resetSeed = useMutation({
    mutationFn: () => api("/dev/reset-seed", { method: "POST" }),
    onSuccess: async () => {
      await qc.invalidateQueries();
      setMsg("Datos de demostración restablecidos.");
    },
    onError: (e: Error) => setMsg(e.message),
  });

  const save = useMutation({
    mutationFn: () =>
      api<Congregation>("/congregation", {
        method: "PUT",
        body: JSON.stringify(form),
      }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["congregation"] });
      setMsg("Congregación guardada.");
    },
    onError: (e: Error) => setMsg(e.message),
  });

  async function changePassword(e: FormEvent) {
    e.preventDefault();
    setMsg("");
    try {
      await api("/auth/change-password", {
        method: "POST",
        body: JSON.stringify(pw),
      });
      setPw({ currentPassword: "", newPassword: "" });
      setMsg("Contraseña actualizada.");
    } catch (err) {
      setMsg(err instanceof Error ? err.message : "No se pudo cambiar la contraseña");
    }
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl">Congregación</h1>
      {msg ? <p className="text-sm text-primary">{msg}</p> : null}
      <Card className="max-w-xl space-y-4">
        <CardTitle>Datos</CardTitle>
        <div className="space-y-1">
          <Label>Nombre</Label>
          <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
        </div>
        <div className="space-y-1">
          <Label>Número</Label>
          <Input value={form.number} onChange={(e) => setForm({ ...form, number: e.target.value })} />
        </div>
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-1">
            <Label>Reunión entre semana</Label>
            <NativeSelect
              className="w-full"
              value={String(form.midweekDay)}
              onChange={(e) => setForm({ ...form, midweekDay: Number(e.target.value) })}
            >
              {weekdayLabels.map((d, i) => (
                <NativeSelectOption key={d} value={i}>
                  {d}
                </NativeSelectOption>
              ))}
            </NativeSelect>
          </div>
          <div className="space-y-1">
            <Label>Reunión de fin de semana</Label>
            <NativeSelect
              className="w-full"
              value={String(form.weekendDay)}
              onChange={(e) => setForm({ ...form, weekendDay: Number(e.target.value) })}
            >
              {weekdayLabels.map((d, i) => (
                <NativeSelectOption key={d} value={i}>
                  {d}
                </NativeSelectOption>
              ))}
            </NativeSelect>
          </div>
        </div>
        <Button onClick={() => save.mutate()} disabled={save.isPending}>
          Guardar
        </Button>
      </Card>

      <Card className="max-w-xl space-y-4">
        <CardTitle>Cambiar contraseña</CardTitle>
        <form className="space-y-3" onSubmit={(e) => void changePassword(e)}>
          <div className="space-y-1">
            <Label>Contraseña actual</Label>
            <Input type="password" value={pw.currentPassword} onChange={(e) => setPw({ ...pw, currentPassword: e.target.value })} required />
          </div>
          <div className="space-y-1">
            <Label>Nueva contraseña</Label>
            <Input type="password" minLength={8} value={pw.newPassword} onChange={(e) => setPw({ ...pw, newPassword: e.target.value })} required />
          </div>
          <Button type="submit">Actualizar</Button>
        </form>
      </Card>

      {q.data?.seedResetEnabled ? (
        <Card className="max-w-xl space-y-3 border-destructive/30">
          <CardTitle>Datos de demostración</CardTitle>
          <p className="text-sm text-muted-foreground">
            Borra publicadores, familias, informes y asistencia y vuelve a cargar el ejemplo. Tu usuario no se toca.
          </p>
          <Button variant="destructive" disabled={resetSeed.isPending} onClick={() => setResetOpen(true)}>
            {resetSeed.isPending ? "Restableciendo…" : "Restablecer datos de demostración"}
          </Button>
          <AlertDialog open={resetOpen} onOpenChange={setResetOpen}>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>¿Restablecer los datos de demostración?</AlertDialogTitle>
                <AlertDialogDescription>
                  Se perderán los publicadores, familias, informes y asistencia que hayas cambiado. Tu usuario no se toca.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancelar</AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  onClick={() => {
                    resetSeed.mutate();
                    setResetOpen(false);
                  }}
                >
                  Restablecer
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </Card>
      ) : null}
    </div>
  );
}
