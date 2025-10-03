import React, { useEffect, useState, useRef } from "react";
import { Card, CardContent } from "@/components/ui/card";
import http from "@/http";
import clsx from "clsx";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { NeoButton } from "@/components/ui/neo-button";
import { deleteJob } from "@/scripts/sheets/delete";

function useDebounce(value: string, delay: number) {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const handler = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(handler);
  }, [value, delay]);
  return debounced;
}

export default function SheetListCard() {
  const [items, setItems] = useState<any[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const debouncedSearch = useDebounce(search, 400);

  // Fetch queue items
  useEffect(() => {
    setLoading(true);
    // Disable cache for search by adding a random param
    const url = `/api/v1/sheets/get?search=${encodeURIComponent(debouncedSearch)}&latest=true&obj_num=5&_nocache=${Date.now()}`;
    http.get(url).then((res) => {
      setItems(res.data || []);
      setLoading(false);
    }).catch(() => setLoading(false));
  }, [debouncedSearch]);

  const handleDelete = async (id: string) => {
    setDeletingId(id);
    const res = await deleteJob(id);
    if (res.status === "deleted") {
      toast.success("Sheet deleted!");
      setItems(items => items.filter(item => item.id !== id));
    } else {
      toast.error(res.error || "Failed to delete sheet.");
    }
    setDeletingId(null);
  };

  return (
    <Card className="max-w-full border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
      <CardContent className="mt-1 px-2 py-2">
        <h1 className="text-5xl font-extrabold tracking-tight mb-4" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
          Sheet List
        </h1>
        <input
          className="mb-4 w-full px-3 py-2 border-2 border-black rounded-lg font-medium text-base"
          placeholder="Search sheets..."
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        {loading ? (
          <div className="flex flex-col items-center justify-center py-8">
            <Loader2 className="h-10 w-10 text-black animate-spin mb-2" />
            <div className="text-lg font-medium text-gray-700">Loading sheets...</div>
          </div>
        ) : items.length === 0 ? (
          <div className="text-base font-medium text-gray-600">No sheets found.</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {items.map((item, idx) => (
              <div
                key={item.id || idx}
                className={clsx(
                  "relative border-2 border-black rounded-lg p-3 bg-gray-50 shadow-[2px_2px_0_0_#000] flex flex-col h-full",
                  item.status === "completed"
                    ? "border-green-600"
                    : item.status === "retrying"
                    ? "border-yellow-600"
                    : item.status === "processing"
                    ? "border-blue-600"
                    : "border-red-600"
                )}
                onClick={() => {
                  if (item.status === "completed" && item.result && item.result.pdf_url) {
                    toast.success("Redirecting to content...");
                    window.open(item.result.pdf_url, "_blank");
                  }
                }}
              >
                {/* Overlay loader for processing */}
                {item.status === "processing" && (
                  <div className="absolute inset-0 z-10 flex flex-col items-center justify-center bg-white/70 backdrop-blur-sm rounded-lg">
                    <Loader2 className="h-10 w-10 text-blue-600 animate-spin mb-2" />
                    <span className="text-lg font-bold text-blue-700">Processing...</span>
                    <NeoButton
                      color="blue"
                      title={"Open Status"}
                      //disabled={deletingId === item.id}
                      onClick={e => {
                      e.stopPropagation();
                      window.location.href = `/sheets/status?jobId=${item.id}`;
                      }}
                      className="text-xs px-4 py-1"
                    />
                    <NeoButton
                      color="red"
                      title={"Delete"}
                      //disabled={deletingId === item.id}
                      onClick={e => {
                      e.stopPropagation();
                      handleDelete(item.id);
                      }}
                      className="text-xs px-4 mt-2 py-1"
                    />
                  </div>
                )}
                {item.status === "retrying" && (
                  <div className="absolute inset-0 z-10 flex flex-col items-center justify-center bg-white/70 backdrop-blur-none rounded-lg">
                    <Loader2 className="h-10 w-10 text-yellow-600 animate-spin mb-2" />
                    <span className="text-lg font-bold text-yellow-700">Retrying...</span>
                    <NeoButton
                      color="blue"
                      title={"Open Status"}
                      //disabled={deletingId === item.id}
                      onClick={e => {
                      e.stopPropagation();
                      window.location.href = `/sheets/status?jobId=${item.id}`;
                      }}
                      className="text-xs px-4 py-1"
                    />
                    
                  </div>
                )}
                {/* Card content below */}
                <div className={item.status === "processing" ? "pointer-events-none opacity-50" : ""}>
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-bold text-lg truncate">
                      {item.prompt ? JSON.parse(item.prompt).subject : "Untitled"}
                    </span>
                    <span
                      className={clsx(
                        "px-2 py-1 rounded font-bold text-xs whitespace-nowrap ml-2",
                        item.status === "completed"
                          ? "bg-green-200 text-green-800"
                          : item.status === "retrying"
                          ? "bg-yellow-200 text-yellow-800"
                          : item.status === "processing"
                          ? "bg-blue-200 text-blue-800"
                          : "bg-red-200 text-red-800"
                      )}
                    >
                      {item.status || "unknown"}
                    </span>
                  </div>
                  <div className="mb-2 text-sm text-gray-700 line-clamp-2">
                    {item.prompt ? JSON.parse(item.prompt).description : ""}
                  </div>
                  <div className="mb-2 text-xs text-gray-500">
                    Created: {item.created_at ? new Date(item.created_at).toLocaleString() : "N/A"}
                  </div>
                  {item.status === "completed" && item.result && item.result.pdf_url && (
                    <div className="mb-2 flex-grow">
                      <iframe
                        src={item.result.pdf_url}
                        title={`PDF-${item.id}`}
                        className="w-full h-40 border border-black rounded"
                      />
                    </div>
                  )}
                  {item.result && typeof item.result === "string" && (
                    <div className="mb-2 text-xs text-red-600 font-mono">{item.result}</div>
                  )}
                  {item.result && item.result.metadata && (
                    <div className="mb-2 bg-gray-100 border border-gray-300 rounded p-2">
                      <div className="font-bold text-xs mb-1">Metadata:</div>
                      <div className="text-xs overflow-hidden">
                        {item.result.metadata.generated && (
                          <div>
                            Generated: {new Date(item.result.metadata.generated).toLocaleString()}
                          </div>
                        )}
                        {item.result.metadata.source && (
                          <div>Source: {item.result.metadata.source}</div>
                        )}
                      </div>
                    </div>
                  )}
                  <div className="flex gap-2 text-xs text-gray-600 mt-auto">
                    <span>Retries: {item.retries ?? 0}/{item.max_retry ?? 0}</span>
                    <span>User: {item.user_id ?? "?"}</span>
                  </div>
                  {/* Delete button for all sheets */}
                    <div className="absolute bottom-3 right-3 z-20 flex gap-2">
                    {item.status !== "retrying" && (
                      <NeoButton
                        color="blue"
                        title="Open"
                        onClick={e => {
                        e.stopPropagation();
                        toast.success("If possible, redirecting to content...");
                        window.open(item.result.pdf_url, "_blank");
                        
                        }}
                        className="text-xs px-4 py-1"
                      />
                    )}
                    <NeoButton
                      color="red"
                      title={deletingId === item.id ? "Deleting..." : "Delete"}
                      disabled={deletingId === item.id}
                      onClick={e => {
                      e.stopPropagation();
                      handleDelete(item.id);
                      }}
                      className="text-xs px-4 py-1"
                    />
                    </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}