import { useMemo, useState, type FormEvent, type ReactNode } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { Card, CardTitle } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { api, type Household, type Publisher } from "@/lib/api";
import { activityLabels, genderLabels, spiritualLabels } from "@/lib/labels";

type PublisherForm = {
  firstName: string;
  lastName: string;
  gender: "male" | "female";
  phone: string;
  email: string;
  householdId: string;
  baptismDate: string;
  startedPreachingDate: string;
  spiritualStatus: "student" | "unbaptized_publisher" | "publisher";
  isElder: boolean;
  isMinisterialServant: boolean;
  isRegularPioneer: boolean;
  isSpecialPioneer: boolean;
};

const emptyPublisher: PublisherForm = {
  firstName: "",
  lastName: "",
  gender: "male",
  phone: "",
  email: "",
  householdId: "",
  baptismDate: "",
  startedPreachingDate: "",
  spiritualStatus: "publisher",
  isElder: false,
  isMinisterialServant: false,
  isRegularPioneer: false,
  isSpecialPioneer: false,
};

function toForm(p: Publisher): PublisherForm {
  return {
    firstName: p.firstName,
    lastName: p.lastName,
    gender: p.gender,
    phone: p.phone,
    email: p.email,
    householdId: p.householdId ?? "",
    baptismDate: p.baptismDate ?? "",
    startedPreachingDate: p.startedPreachingDate ?? "",
    spiritualStatus: p.spiritualStatus,
    isElder: p.isElder,
    isMinisterialServant: p.isMinisterialServant,
    isRegularPioneer: p.isRegularPioneer,
    isSpecialPioneer: p.isSpecialPioneer,
  };
}

function payload(f: PublisherForm) {
  return {
    firstName: f.firstName,
    lastName: f.lastName,
    gender: f.gender,
    phone: f.phone,
    email: f.email,
    householdId: f.householdId || null,
    baptismDate: f.baptismDate || null,
    startedPreachingDate: f.startedPreachingDate || null,
    spiritualStatus: f.spiritualStatus,
    isElder: f.isElder,
    isMinisterialServant: f.isMinisterialServant,
    isRegularPioneer: f.isRegularPioneer,
    isSpecialPioneer: f.isSpecialPioneer,
  };
}

const col = createColumnHelper<Publisher>();

export function PublishersPage() {
  const qc = useQueryClient();
  const [tab, setTab] = useState<"people" | "homes">("people");
  const [q, setQ] = useState("");
  const [activity, setActivity] = useState("");
  const [privilege, setPrivilege] = useState("");
  const [editing, setEditing] = useState<Publisher | null | "new">(null);
  const [form, setForm] = useState<PublisherForm>(emptyPublisher);

  const pubs = useQuery({ queryKey: ["publishers"], queryFn: () => api<Publisher[]>("/publishers") });
  const homes = useQuery({ queryKey: ["households"], queryFn: () => api<Household[]>("/households") });

  const filtered = useMemo(() => {
    return (pubs.data ?? []).filter((p) => {
      const hay = `${p.firstName} ${p.lastName} ${p.householdName ?? ""}`.toLowerCase();
      if (q && !hay.includes(q.toLowerCase())) return false;
      if (activity === "active" && !p.isActive) return false;
      if (activity && activity !== "active" && p.activityStatus !== activity) return false;
      if (privilege === "elder" && !p.isElder) return false;
      if (privilege === "ms" && !p.isMinisterialServant) return false;
      if (privilege === "rp" && !p.isRegularPioneer) return false;
      if (privilege === "sp" && !p.isSpecialPioneer) return false;
      return true;
    });
  }, [pubs.data, q, activity, privilege]);

  const columns = useMemo(
    () => [
      col.accessor((r) => `${r.lastName}, ${r.firstName}`, { id: "name", header: "Nombre" }),
      col.accessor((r) => genderLabels[r.gender], { id: "gender", header: "Sexo" }),
      col.accessor((r) => r.householdName ?? "—", { id: "home", header: "Familia" }),
      col.accessor((r) => spiritualLabels[r.spiritualStatus], { id: "spirit", header: "Estado" }),
      col.accessor((r) => activityLabels[r.activityStatus], { id: "act", header: "Actividad" }),
      col.display({
        id: "priv",
        header: "Privilegios",
        cell: ({ row }) => (
          <div className="flex flex-wrap gap-1">
            {row.original.isElder ? <Badge>Anciano</Badge> : null}
            {row.original.isMinisterialServant ? <Badge>SM</Badge> : null}
            {row.original.isRegularPioneer ? <Badge>PR</Badge> : null}
            {row.original.isSpecialPioneer ? <Badge>PE</Badge> : null}
          </div>
        ),
      }),
      col.display({
        id: "actions",
        header: "",
        cell: ({ row }) => (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              setEditing(row.original);
              setForm(toForm(row.original));
            }}
          >
            Editar
          </Button>
        ),
      }),
    ],
    [],
  );

  const table = useReactTable({
    data: filtered,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  const savePub = useMutation({
    mutationFn: async () => {
      if (editing && editing !== "new") {
        return api(`/publishers/${editing.id}`, { method: "PUT", body: JSON.stringify(payload(form)) });
      }
      return api("/publishers", { method: "POST", body: JSON.stringify(payload(form)) });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["publishers"] });
      setEditing(null);
    },
  });

  const delPub = useMutation({
    mutationFn: (id: string) => api(`/publishers/${id}`, { method: "DELETE" }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["publishers"] });
      setEditing(null);
    },
  });

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="text-3xl">Publicadores</h1>
        <div className="flex gap-2">
          <Button variant={tab === "people" ? "default" : "outline"} onClick={() => setTab("people")}>
            Personas
          </Button>
          <Button variant={tab === "homes" ? "default" : "outline"} onClick={() => setTab("homes")}>
            Familias
          </Button>
        </div>
      </header>

      {tab === "homes" ? (
        <HouseholdsPanel homes={homes.data ?? []} />
      ) : (
        <>
          <div className="flex flex-wrap gap-2">
            <Input placeholder="Buscar…" value={q} onChange={(e) => setQ(e.target.value)} className="max-w-xs" />
            <NativeSelect value={activity} onChange={(e) => setActivity(e.target.value)}>
              <NativeSelectOption value="">Actividad: todas</NativeSelectOption>
              <NativeSelectOption value="active">Activos</NativeSelectOption>
              <NativeSelectOption value="regular">Regulares</NativeSelectOption>
              <NativeSelectOption value="irregular">Irregulares</NativeSelectOption>
              <NativeSelectOption value="inactive">Inactivos</NativeSelectOption>
            </NativeSelect>
            <NativeSelect value={privilege} onChange={(e) => setPrivilege(e.target.value)}>
              <NativeSelectOption value="">Privilegios: todos</NativeSelectOption>
              <NativeSelectOption value="elder">Ancianos</NativeSelectOption>
              <NativeSelectOption value="ms">Siervos ministeriales</NativeSelectOption>
              <NativeSelectOption value="rp">Precursores regulares</NativeSelectOption>
              <NativeSelectOption value="sp">Precursores especiales</NativeSelectOption>
            </NativeSelect>
            <Button
              onClick={() => {
                setEditing("new");
                setForm(emptyPublisher);
              }}
            >
              Nuevo publicador
            </Button>
          </div>

          <div className="overflow-x-auto rounded-lg border border-border bg-card">
            <table className="w-full text-left text-sm">
              <thead className="bg-muted">
                {table.getHeaderGroups().map((hg) => (
                  <tr key={hg.id}>
                    {hg.headers.map((h) => (
                      <th key={h.id} className="px-3 py-2 font-medium">
                        {flexRender(h.column.columnDef.header, h.getContext())}
                      </th>
                    ))}
                  </tr>
                ))}
              </thead>
              <tbody>
                {table.getRowModel().rows.map((row) => (
                  <tr key={row.id} className="border-t border-border">
                    {row.getVisibleCells().map((cell) => (
                      <td key={cell.id} className="px-3 py-2">
                        {flexRender(cell.column.columnDef.cell, cell.getContext())}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}

      {editing ? (
        <div className="fixed inset-0 z-20 flex items-start justify-center overflow-auto bg-foreground/40 p-4">
          <Card className="mt-10 w-full max-w-xl space-y-4">
            <CardTitle>{editing === "new" ? "Nuevo publicador" : "Editar publicador"}</CardTitle>
            <form
              className="grid gap-3 sm:grid-cols-2"
              onSubmit={(e) => {
                e.preventDefault();
                savePub.mutate();
              }}
            >
              <Field label="Nombre">
                <Input required value={form.firstName} onChange={(e) => setForm({ ...form, firstName: e.target.value })} />
              </Field>
              <Field label="Apellidos">
                <Input required value={form.lastName} onChange={(e) => setForm({ ...form, lastName: e.target.value })} />
              </Field>
              <Field label="Sexo">
                <NativeSelect className="w-full" value={form.gender} onChange={(e) => setForm({ ...form, gender: e.target.value as PublisherForm["gender"] })}>
                  <NativeSelectOption value="male">Varón</NativeSelectOption>
                  <NativeSelectOption value="female">Mujer</NativeSelectOption>
                </NativeSelect>
              </Field>
              <Field label="Familia">
                <NativeSelect className="w-full" value={form.householdId} onChange={(e) => setForm({ ...form, householdId: e.target.value })}>
                  <NativeSelectOption value="">Sin familia</NativeSelectOption>
                  {(homes.data ?? []).map((h) => (
                    <NativeSelectOption key={h.id} value={h.id}>
                      {h.name}
                    </NativeSelectOption>
                  ))}
                </NativeSelect>
              </Field>
              <Field label="Teléfono">
                <Input value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
              </Field>
              <Field label="Correo">
                <Input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
              </Field>
              <Field label="Bautismo">
                <Input type="date" value={form.baptismDate} onChange={(e) => setForm({ ...form, baptismDate: e.target.value })} />
              </Field>
              <Field label="Empezó a predicar">
                <Input type="date" value={form.startedPreachingDate} onChange={(e) => setForm({ ...form, startedPreachingDate: e.target.value })} />
              </Field>
              <Field label="Estado espiritual" className="sm:col-span-2">
                <NativeSelect
                  className="w-full"
                  value={form.spiritualStatus}
                  onChange={(e) => setForm({ ...form, spiritualStatus: e.target.value as PublisherForm["spiritualStatus"] })}
                >
                  <NativeSelectOption value="publisher">Publicador</NativeSelectOption>
                  <NativeSelectOption value="unbaptized_publisher">Publicador no bautizado</NativeSelectOption>
                  <NativeSelectOption value="student">Estudiante</NativeSelectOption>
                </NativeSelect>
              </Field>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox checked={form.isElder} onCheckedChange={(checked) => setForm({ ...form, isElder: checked === true, isMinisterialServant: checked === true ? false : form.isMinisterialServant })} />
                Anciano
              </label>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox checked={form.isMinisterialServant} onCheckedChange={(checked) => setForm({ ...form, isMinisterialServant: checked === true, isElder: checked === true ? false : form.isElder })} />
                Siervo ministerial
              </label>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox checked={form.isRegularPioneer} onCheckedChange={(checked) => setForm({ ...form, isRegularPioneer: checked === true })} />
                Precursor regular
              </label>
              <label className="flex items-center gap-2 text-sm">
                <Checkbox checked={form.isSpecialPioneer} onCheckedChange={(checked) => setForm({ ...form, isSpecialPioneer: checked === true })} />
                Precursor especial
              </label>
              {savePub.isError ? <p className="sm:col-span-2 text-sm text-destructive">{savePub.error.message}</p> : null}
              <div className="sm:col-span-2 flex justify-between">
                {editing !== "new" ? (
                  <Button variant="destructive" onClick={() => delPub.mutate(editing.id)}>
                    Eliminar
                  </Button>
                ) : (
                  <span />
                )}
                <div className="flex gap-2">
                  <Button variant="outline" onClick={() => setEditing(null)}>
                    Cancelar
                  </Button>
                  <Button type="submit" disabled={savePub.isPending}>
                    Guardar
                  </Button>
                </div>
              </div>
            </form>
          </Card>
        </div>
      ) : null}
    </div>
  );
}

function Field({ label, children, className }: { label: string; children: ReactNode; className?: string }) {
  return (
    <div className={className}>
      <Label>{label}</Label>
      <div className="mt-1">{children}</div>
    </div>
  );
}

function HouseholdsPanel({ homes }: { homes: Household[] }) {
  const qc = useQueryClient();
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [notes, setNotes] = useState("");
  const [editing, setEditing] = useState<Household | null>(null);

  const save = useMutation({
    mutationFn: async () => {
      const body = JSON.stringify({ name, address, notes });
      if (editing) return api(`/households/${editing.id}`, { method: "PUT", body });
      return api("/households", { method: "POST", body });
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["households"] });
      setEditing(null);
      setName("");
      setAddress("");
      setNotes("");
    },
  });

  const del = useMutation({
    mutationFn: (id: string) => api(`/households/${id}`, { method: "DELETE" }),
    onSuccess: async () => qc.invalidateQueries({ queryKey: ["households"] }),
  });

  function startEdit(h: Household) {
    setEditing(h);
    setName(h.name);
    setAddress(h.address);
    setNotes(h.notes);
  }

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    save.mutate();
  }

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
      <div className="space-y-2">
        {homes.map((h) => (
          <Card key={h.id} className="flex items-start justify-between gap-3">
            <div>
              <div className="font-medium">{h.name}</div>
              <div className="text-sm text-muted-foreground">{h.address || "Sin dirección"}</div>
              {h.notes ? <div className="mt-1 text-sm">{h.notes}</div> : null}
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={() => startEdit(h)}>
                Editar
              </Button>
              <Button variant="ghost" size="sm" onClick={() => del.mutate(h.id)}>
                Borrar
              </Button>
            </div>
          </Card>
        ))}
      </div>
      <Card>
        <CardTitle>{editing ? "Editar familia" : "Nueva familia"}</CardTitle>
        <form className="mt-4 space-y-3" onSubmit={onSubmit}>
          <div>
            <Label>Nombre</Label>
            <Input className="mt-1" required value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div>
            <Label>Dirección</Label>
            <Input className="mt-1" value={address} onChange={(e) => setAddress(e.target.value)} />
          </div>
          <div>
            <Label>Notas</Label>
            <Textarea className="mt-1" value={notes} onChange={(e) => setNotes(e.target.value)} />
          </div>
          <div className="flex gap-2">
            <Button type="submit">{editing ? "Guardar" : "Crear"}</Button>
            {editing ? (
              <Button
                variant="outline"
                onClick={() => {
                  setEditing(null);
                  setName("");
                  setAddress("");
                  setNotes("");
                }}
              >
                Cancelar
              </Button>
            ) : null}
          </div>
        </form>
      </Card>
    </div>
  );
}
