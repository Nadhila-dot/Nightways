export function connectToJobWebSocket(jobId: string, sessionId: string, onUpdate: (data: any) => void) {
  const wsUrl = `${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/api/v1/ws/job/${jobId}?session=${sessionId}`;
  const ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log(`Connected to WebSocket for job ${jobId}`);
  };

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      console.log("WebSocket message received:", data);
      onUpdate(data);
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  };

  ws.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  ws.onclose = () => {
    console.log(`WebSocket connection for job ${jobId} closed`);
  };

  return ws; // Return the WebSocket instance for cleanup if needed
}