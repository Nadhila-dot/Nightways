import React from "react";
import { Card, CardHeader, CardContent } from "@/components/ui/card";

export default function RefreshCache() {
  function handleCache() {
        localStorage.removeItem("http_cache");
        window.location.reload(); // rebuild new cache from api.
  }

  return (
    <Card className="max-w-3xl border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
      
      <CardContent className="mt-1 px-2 py-2">
        <h1 className="text-5xl font-extrabold tracking-tight" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
        Refresh Cache
      </h1>
        <div className="mb-4 text-base font-medium">Click below to refresh the cache. This will rebuild the cache so you can get the latest data.</div>
        <button
          className="bg-yellow-600 text-white border-2 border-black rounded-lg px-6 py-2 font-bold shadow-[4px_4px_0_0_#000] hover:bg-yellow-700 transition-all"
          onClick={handleCache}
        >
          Refresh
        </button>
      </CardContent>
    </Card>
  );
}