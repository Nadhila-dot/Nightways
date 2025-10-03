import http from "@/http";

export async function addSheetToNotebook(notebookId: number, sheetName: string, url: string) {
  // You may need to pass the notebookId and sheet info
  const res = await http.post(`/api/v1/notebooks/${notebookId}/items`, {
    sheetName,
    url,
  });
  return res.data;
}