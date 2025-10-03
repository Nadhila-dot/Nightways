import { useEffect } from "react";
import { toast } from "sonner";
import { Info, AlertCircle, CheckCircle, Settings } from "lucide-react";

function getIcon(type: string) {
  switch (type) {
    case "error":
      return <AlertCircle className="text-red-600" />;
    case "info":
      return <Info className="text-blue-600" />;
    case "success":
      return <CheckCircle className="text-green-600" />;
    case "settings":
      return <Settings className="text-yellow-600" />;
    default:
      return <Info className="text-gray-600" />;
  }
}

export function NotificationListener({ sessionId }: { sessionId: string }) {
  useEffect(() => {
    if (!sessionId) return;

    const wsUrl = `${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/api/v1/ws/notifications?session=${sessionId}`;
    console.log("Connecting to WebSocket at:", wsUrl);
    const ws = new WebSocket(wsUrl);

    ws.onmessage = (event) => {
      try {
        console.log("WebSocket message received:", event.data);
        const msg = JSON.parse(event.data);

        // Unified handler for push/broadcast/notification
        let payload = msg.data?.data || msg.data || {};
        let type = payload.type || msg.type || "info";
        let message = payload.message || msg.message || "Notification";
        let timestamp = msg.timestamp || new Date().toISOString();

        toast(message, {
          description: `${type.charAt(0).toUpperCase() + type.slice(1)} • ${new Date(timestamp).toLocaleTimeString()}`,
          icon: getIcon(type),
          duration: 4000,
        });
      } catch (e) {
        console.error("Error processing websocket message:", e);
      }
    };

    ws.onopen = () => {
      console.log("WebSocket connection established");
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      toast("WebSocket error", { description: "Notifications unavailable", duration: 300 });
    };

    ws.onclose = (event) => {
      // Only show toast for abnormal closes
      if (event.code !== 1000 && event.code !== 1001) {
      //  toast("WebSocket closed", { description: "Notifications disconnected", duration: 3000 });
      console.log("WebSocket closed");
      }
      console.log("WebSocket connection closed", event);
    };

    return () => ws.close();
  }, [sessionId]);

  return null;
}