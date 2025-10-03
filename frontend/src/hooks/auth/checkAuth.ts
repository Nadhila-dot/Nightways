import { useState, useEffect } from "react";
import { toast } from "sonner";
import { QueryClient } from "@tanstack/query-core";
import http from "@/http";
import { time } from "console";


const queryClient = new QueryClient();

async function fetchSessionStatus(session: string) {
  const res = await http.get(`/api/v1/session/${session}`);
  const time = Date.now();
  localStorage.setItem(
    "auth-data",
    JSON.stringify({ ...res.data, time })
  ); 
  return res.status === 200;

}

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const session = localStorage.getItem("session");
    const path = window.location.pathname;
    const isAuthRoute = path.startsWith("/auth");
    const isIndexRoute = path.startsWith("/index");

    if (!session) {
      setIsAuthenticated(false);
      if (!isAuthRoute && !isIndexRoute) {
        window.location.href = "/auth/login";
      }
      return;
    }

    queryClient
      .fetchQuery({
        queryKey: ["sessionStatus", session],
        queryFn: () => fetchSessionStatus(session),
        staleTime: 1000 * 60, // 1 minute
      })
      .then((valid) => {
        setIsAuthenticated(valid);
        if (valid) {
        //  toast.success("Session validated");
        } else if (!isAuthRoute && !isIndexRoute) {
          window.location.href = "/auth/login";
        }
      })
      .catch(() => {
        setIsAuthenticated(false);
        if (!isAuthRoute && !isIndexRoute) {
          window.location.href = "/auth/login";
        }
      });
  }, []);

  return isAuthenticated;
}