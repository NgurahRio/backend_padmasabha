export const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function api(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers);
  if (!(init.body instanceof FormData)) headers.set("Content-Type", "application/json");
  const response = await fetch(`${API_URL}${path}`, { ...init, headers, credentials: "include" });
  const text = await response.text();
  let payload: any = null;
  try { payload = text ? JSON.parse(text) : null; } catch { payload = text; }
  if (!response.ok) {
    if (response.status === 401 && typeof window !== "undefined" && path !== "/admin/login") window.dispatchEvent(new Event("admin-session-expired"));
    throw new Error(payload?.message || payload?.error || `Request failed (${response.status})`);
  }
  return payload;
}

export const unwrap = (payload: any) => Array.isArray(payload) ? payload : (payload?.data ?? []);
export const idOf = (row: any) => row?.id_destination ?? row?.id_event ?? row?.id_categories ?? row?.id_users ?? row?.id_sos ?? row?.id_facility ?? row?.id_subpackage ?? row?.id_packages;
