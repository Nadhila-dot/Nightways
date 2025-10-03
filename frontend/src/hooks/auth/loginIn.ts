import http from "@/http";

export async function loginIn(identifier: string, password: string) {
  const res = await http.post("/api/v1/login", { identifier, password });

  const data = res.data;

  if (!res.status || res.status < 200 || res.status >= 300) {
    throw new Error(data.error || "Login failed");
  }

  // Store session key for future requests
  if (data.session) {
    localStorage.setItem("session", data.session);
  }

  // Returns session ID
  return {
    session: data.session
  };
}