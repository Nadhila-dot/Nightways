import { BrowserRouter as Router, Routes, Route, Link, useLocation, RouterProvider } from "react-router-dom";
import CountBtn from "@/components/count-btn";
import ReactSVG from "@/assets/react.svg";
import { Badge } from "@/components/ui/badge";
import { NotFound } from "./components/ScreenBlock";
import { HomeContainer } from "./containers/index/container";
import { LoginContainer } from "./containers/login/loginContainer";
import { Suspense, useEffect } from "react";
import JsonResponse from "./json/json-route";
import { RouterIcon } from "lucide-react";

import { createBrowserRouter} from "react-router-dom";
import Sidebar from "./components/Sidebar/Sidebar";
import { useAuth } from "./hooks/auth/checkAuth";

import InvalidSessionModal from "./containers/session/Invalid";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useSearch } from "./hooks/search/useSearch";
import { useSystemStore } from './state/system';
import { Toaster } from "./components/ui/sonner";
import { NotificationListener } from "./components/Notification/handle";


function App() {
  const isAuthenticated = useAuth();
  const queryClient = new QueryClient();
  const sessionId = localStorage.getItem("session") || ""; // Or get from your auth system

  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <Toaster position="bottom-right" />
        <NotificationListener sessionId={sessionId} />
        {!(window.location.pathname.startsWith("/auth") || window.location.pathname.startsWith("/index")) && (
          <Sidebar isAuthenticated={isAuthenticated} />
        )}
        {!isAuthenticated &&
          !window.location.pathname.startsWith("/auth") &&
          !window.location.pathname.startsWith("/index") && (
            <InvalidSessionModal />
        )}
        <Routes>
          <Route path="/auth/login" element={<LoginContainer />} />
        </Routes>
        
      </Router>
    </QueryClientProvider>
  );
}

export default App;