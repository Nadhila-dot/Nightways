import { useState, useEffect } from "react";

export function useSession() {
  const [session, setSession] = useState<string | null>(null);

  useEffect(() => {
    const storedSession = localStorage.getItem("session");
    setSession(storedSession);
  }, []);

  return session;
}