import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  Outlet,
  RouterProvider,
  createRootRoute,
  createRoute,
  createRouter,
  redirect,
} from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { AppLayout } from "@/components/layout";
import { AttendancePage } from "@/features/attendance";
import { CongregationPage } from "@/features/congregation";
import { DashboardPage } from "@/features/dashboard";
import { LoginPage } from "@/features/login";
import { PublishersPage } from "@/features/publishers";
import { ReportsPage } from "@/features/reports";
import { api, ApiError } from "@/lib/api";
import "./index.css";

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: false, refetchOnWindowFocus: false } },
});

async function requireSession() {
  try {
    await api("/api/auth/me");
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      throw redirect({ to: "/login" });
    }
    throw err;
  }
}

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: LoginPage,
});

const appRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: "app",
  beforeLoad: requireSession,
  component: AppLayout,
});

const dashboardRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/",
  component: DashboardPage,
});

const publishersRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/publishers",
  component: PublishersPage,
});

const reportsRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/reports",
  component: ReportsPage,
});

const attendanceRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/attendance",
  component: AttendancePage,
});

const congregationRoute = createRoute({
  getParentRoute: () => appRoute,
  path: "/congregation",
  component: CongregationPage,
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  appRoute.addChildren([dashboardRoute, publishersRoute, reportsRoute, attendanceRoute, congregationRoute]),
]);

const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </StrictMode>,
);
