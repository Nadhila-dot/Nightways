import http from "@/http";

export async function registerIn(
  username: string,
  email: string,
  password: string,
  rank: string = "user"
) {
  const res = await http.post("/api/v1/register", {
    username,
    email,
    password,
    rank,
  });

  const data = res.data;

  if (!res.status || res.status < 200 || res.status >= 300) {
    throw new Error(data.error || "Registration failed");
  }

  // Returns success message
  return {
    message: data.message,
  };
}