import React from "react";
import { Card, CardHeader, CardContent } from "@/components/ui/card";

export default function LogoutCard() {
  function handleLogout() {
    localStorage.removeItem("session");
    // remove all forms of caches from the user aswell
    // So we can't mess up
    localStorage.removeItem("auth-data");
    localStorage.removeItem("http_cache");
    window.location.reload();
  }

  return (
    <Card className="max-w-3xl border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
      
      <CardContent className="mt-1 px-2 py-2">
        <h1 className="text-5xl font-extrabold tracking-tight" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
        Logout
      </h1>
        <div className="mb-4 text-base font-medium">Click below to log out of your session. All users can have a max of 3 sessions for server load reasons. All cache will also be removed.</div>
        <button
          className="bg-red-600 text-white border-2 border-black rounded-lg px-6 py-2 font-bold shadow-[4px_4px_0_0_#000] hover:bg-red-700 transition-all"
          onClick={handleLogout}
        >
          Logout
        </button>
      </CardContent>
    </Card>
  );
}