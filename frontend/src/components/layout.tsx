import { Link, Outlet, useNavigate, useRouterState } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { CalendarDays, ClipboardList, Home, LogOut, Settings, Users } from "lucide-react";
import { api, type Congregation, type User } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const nav = [
  { to: "/", label: "Inicio", icon: Home },
  { to: "/publishers", label: "Publicadores", icon: Users },
  { to: "/reports", label: "Informes", icon: ClipboardList },
  { to: "/attendance", label: "Asistencia", icon: CalendarDays },
  { to: "/congregation", label: "Congregación", icon: Settings },
];

export function AppLayout() {
  const navigate = useNavigate();
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const cong = useQuery({
    queryKey: ["congregation"],
    queryFn: () => api<Congregation>("/api/congregation"),
  });
  const me = useQuery({
    queryKey: ["me"],
    queryFn: () => api<User>("/api/auth/me"),
  });

  async function logout() {
    await api("/api/auth/logout", { method: "POST" });
    await navigate({ to: "/login" });
  }

  return (
    <div className="min-h-screen md:grid md:grid-cols-[240px_1fr]">
      <aside className="border-b border-border bg-card md:border-b-0 md:border-r">
        <div className="px-5 py-6">
          <div className="font-sans text-2xl font-semibold text-primary">khhub</div>
          <div className="mt-1 text-sm text-muted-foreground">
            {cong.data?.name || "Congregación"}
            {cong.data?.number ? ` · ${cong.data.number}` : ""}
          </div>
        </div>
        <nav className="flex gap-1 overflow-x-auto px-3 pb-3 md:flex-col">
          {nav.map((item) => {
            const active = pathname === item.to;
            const Icon = item.icon;
            return (
              <Link
                key={item.to}
                to={item.to}
                className={cn(
                  "flex items-center gap-2 rounded-md px-3 py-2 text-sm whitespace-nowrap",
                  active ? "bg-primary text-primary-foreground" : "hover:bg-muted",
                )}
              >
                <Icon className="h-4 w-4" />
                {item.label}
              </Link>
            );
          })}
        </nav>
        <div className="hidden items-center justify-between px-4 py-4 md:flex">
          <span className="truncate text-xs text-muted-foreground">{me.data?.email}</span>
          <Button variant="ghost" size="sm" onClick={() => void logout()}>
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </aside>
      <main className="p-4 md:p-8">
        <Outlet />
      </main>
    </div>
  );
}
