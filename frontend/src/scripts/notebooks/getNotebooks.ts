import http from "@/http";

export async function getNotebooks() {
  const { data } = await http.get("/api/v1/notebooks");
  return data;
}