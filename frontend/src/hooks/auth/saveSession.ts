export function saveSession(session: string) {
  if (!session) {
    throw new Error("Invalid session data");
  }
  localStorage.setItem("session", session);
  return true;
}