import React, { useEffect, useState } from "react";
import { Card, CardHeader, CardContent } from "@/components/ui/card";
import { getSystemInfo } from "@/scripts/getSet";

export default function SetCard() {
  const [setData, setSetData] = useState<any>(null);

  useEffect(() => {
    getSystemInfo().then((data) => setSetData(data)); // Setdata is the set data thingy output, it might be confusing
  }, []);

  // Helper to render AI config
  function renderAI(ai: any) {
    if (!ai) return null;
    return (
      <div className="mb-4">
        <div className="font-bold mb-1">AI Configuration:</div>
        <ul className="ml-4 text-base">
          <li><span className="font-semibold">Source:</span> {ai.source}</li>
          <li><span className="font-semibold">Model:</span> {ai.model}</li>
          <li><span className="font-semibold">Temperature:</span> {ai.temperature}</li>
          <li><span className="font-semibold">Max Tokens:</span> {ai.max_tokens}</li>
          <li><span className="font-semibold">System Prompt:</span> <span className="italic">{ai.system_prompt}</span></li>
          <li><span className="font-semibold">Cooldown (sec):</span> {ai.cooldown_sec}</li>
        </ul>
      </div>
    );
  }

  return (
    <Card className="max-w-3xl border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
      
      <CardContent className="mt-1 px-2 py-2">
        <h1 className="text-5xl font-extrabold tracking-tight" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
        System Enviroment (Set.json)
      </h1>
        {!setData ? (
          <div className="text-base font-medium">Loading...</div>
        ) : (
          Array.isArray(setData) && setData.length > 0 ? (
            <div className="p-4">
              <div className="mb-2">
                <span className="font-bold">API Key:</span> <span className="bg-gray-200 px-2 py-1 rounded">{setData[0].AI_API}</span>
              </div>
              {renderAI(setData[0].AI)}
              <div className="mb-2">
                <span className="font-bold">Max Sessions:</span> {setData[0].MAX_SESSIONS}
              </div>
            </div>
          ) : (
            <div className="text-base font-medium">No environment data found.</div>
          )
        )}
      </CardContent>
    </Card>
  );
}