import http from "@/http";

export interface SheetCreateData {
  subject: string;
  course: string;
  description: string;
  tags: string;
  curriculum: string;
  specialInstructions: string;
  visibility: string;
}

export async function createSheet(data: SheetCreateData) {
  const res = await http.post("/api/v1/sheets/create", data);
  return res.data;
}

export async function getSheetQueue() {
  const res = await http.get("/api/v1/sheets/queue");
  return res.data;
}