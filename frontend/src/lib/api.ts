export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers);
  if (init?.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  const res = await fetch(path, {
    ...init,
    credentials: "include",
    headers,
  });
  if (res.status === 204) {
    return undefined as T;
  }
  const data = (await res.json().catch(() => ({}))) as { error?: string };
  if (!res.ok) {
    throw new ApiError(res.status, data.error ?? "Error inesperado");
  }
  return data as T;
}

export type User = { id: string; email: string };

export type Congregation = {
  name: string;
  number: string;
  midweekDay: number;
  weekendDay: number;
  updatedAt?: string;
  seedResetEnabled?: boolean;
};

export type Household = {
  id: string;
  name: string;
  address: string;
  notes: string;
};

export type Publisher = {
  id: string;
  householdId: string | null;
  householdName: string | null;
  firstName: string;
  lastName: string;
  gender: "male" | "female";
  phone: string;
  email: string;
  baptismDate: string | null;
  startedPreachingDate: string | null;
  spiritualStatus: "student" | "unbaptized_publisher" | "publisher";
  isElder: boolean;
  isMinisterialServant: boolean;
  isRegularPioneer: boolean;
  isSpecialPioneer: boolean;
  activityStatus: "regular" | "irregular" | "inactive";
  isActive: boolean;
};

export type ReportRow = {
  publisherId: string;
  firstName: string;
  lastName: string;
  spiritualStatus: string;
  isRegularPioneer: boolean;
  isSpecialPioneer: boolean;
  hourReporter: boolean;
  hasReport: boolean;
  id: string | null;
  sharedInMinistry: boolean;
  bibleStudies: number;
  hours: number | null;
  auxiliaryPioneer: boolean;
  late: boolean;
  remarks: string;
};

export type ReportsResponse = {
  year: number;
  month: number;
  serviceYear: number;
  reports: ReportRow[];
  missing: { publisherId: string; firstName: string; lastName: string }[];
};

export type Attendance = {
  id: string;
  date: string;
  kind: "midweek" | "weekend";
  inPerson: number;
  online: number | null;
};

export type Dashboard = {
  congregation: { name: string; number: string };
  year: number;
  month: number;
  serviceYear: number;
  publishers: {
    total: number;
    students: number;
    active: number;
    regular: number;
    irregular: number;
    inactive: number;
  };
  reports: {
    shouldReport: number;
    reported: number;
    missingCount: number;
    missing: { publisherId: string; firstName: string; lastName: string }[];
    shared: number;
    participation: number;
    bibleStudies: number;
    pioneerHours: number;
  };
  attendance: {
    midweekMeetings: number;
    weekendMeetings: number;
    midweekAvg: number | null;
    weekendAvg: number | null;
    midweekOnlineAvg: number | null;
    weekendOnlineAvg: number | null;
  };
};
