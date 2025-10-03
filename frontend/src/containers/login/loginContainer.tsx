import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { loginIn } from "@/hooks/auth/loginIn";
import { saveSession } from "@/hooks/auth/saveSession";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/auth/checkAuth";
import { registerIn } from "@/hooks/auth/RegisterIn";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

export function LoginContainer() {
  const [tab, setTab] = useState<"login" | "register">("login");
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState<string | null>(null);
  const [msgType, setMsgType] = useState<"error" | "success">("error");

  const isAuthenticated = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      navigate("/home");
      window.location.reload();
    }
  }, [isAuthenticated, navigate]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMsg(null);
    try {
      const res = await loginIn(identifier, password);
      saveSession(res.session);
      setMsgType("success");
      setMsg("Login successful! Redirecting...");
      navigate("/home");
      window.location.reload();
    } catch (err: any) {
      setMsgType("error");
      setMsg(err.message);
    }
    setLoading(false);
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setMsg(null);
    try {
      const res = await registerIn(username, email, password);
      setMsgType("success");
      setMsg("Registration successful! Switching to login...");
      
      // Auto-fill login form and switch to login tab
      setTimeout(() => {
        setIdentifier(username);
        // Keep the password filled
        setTab("login");
        setMsg("You can now login with your credentials!");
      }, 1500);
    } catch (err: any) {
      setMsgType("error");
      setMsg(err.message);
    }
    setLoading(false);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 p-4">
      <Card className="max-w-md w-full border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] bg-white text-black p-0">
        <CardHeader className="px-6 pt-6">
          <div className="flex gap-3 mb-4">
            <button
              className={`flex-1 px-4 py-2 font-bold border-2 border-black rounded-lg transition-all
                ${tab === "login" 
                  ? "bg-black text-white" 
                  : "bg-white text-black hover:bg-gray-100"
                }`}
              onClick={() => setTab("login")}
            >
              Login
            </button>
            <button
              className={`flex-1 px-4 py-2 font-bold border-2 border-black rounded-lg transition-all
                ${tab === "register" 
                  ? "bg-black text-white" 
                  : "bg-white text-black hover:bg-gray-100"
                }`}
              onClick={() => setTab("register")}
            >
              Register
            </button>
          </div>
          <h1 className="text-5xl font-extrabold tracking-tight" style={{ fontFamily: "'Space Grotesk', sans-serif" }}>
            {tab === "login" ? "Welcome Back" : "Join Us"}
          </h1>
        </CardHeader>

        <CardContent className="px-6 pb-6">
          <div className="mb-4 text-base font-medium">
            {tab === "login" 
              ? "Enter your credentials to access your account." 
              : "Create a new account to get started."
            }
          </div>

          {tab === "login" ? (
            <form onSubmit={handleLogin} className="space-y-4">
              <div>
                <label className="block text-sm font-bold mb-2">
                  Username or Email
                </label>
                <input
                  type="text"
                  placeholder="Enter username or email"
                  value={identifier}
                  onChange={e => setIdentifier(e.target.value)}
                  className="w-full border-2 border-black rounded-lg px-4 py-2 font-medium
                    focus:outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-bold mb-2">
                  Password
                </label>
                <input
                  type="password"
                  placeholder="Enter your password"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  className="w-full border-2 border-black rounded-lg px-4 py-2 font-medium
                    focus:outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-blue-600 text-white border-2 border-black rounded-lg px-6 py-2 font-bold 
                  shadow-[4px_4px_0_0_#000] hover:bg-blue-700 transition-all
                  active:shadow-[2px_2px_0_0_#000] active:translate-x-[2px] active:translate-y-[2px]
                  disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? "Logging in..." : "Login"}
              </button>
            </form>
          ) : (
            <form onSubmit={handleRegister} className="space-y-4">
              <div>
                <label className="block text-sm font-bold mb-2">
                  Username
                </label>
                <input
                  type="text"
                  placeholder="Choose a username"
                  value={username}
                  onChange={e => setUsername(e.target.value)}
                  className="w-full border-2 border-black rounded-lg px-4 py-2 font-medium
                    focus:outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-bold mb-2">
                  Email
                </label>
                <input
                  type="email"
                  placeholder="Enter your email"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  className="w-full border-2 border-black rounded-lg px-4 py-2 font-medium
                    focus:outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-bold mb-2">
                  Password
                </label>
                <input
                  type="password"
                  placeholder="Create a password"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  className="w-full border-2 border-black rounded-lg px-4 py-2 font-medium
                    focus:outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-green-600 text-white border-2 border-black rounded-lg px-6 py-2 font-bold 
                  shadow-[4px_4px_0_0_#000] hover:bg-green-700 transition-all
                  active:shadow-[2px_2px_0_0_#000] active:translate-x-[2px] active:translate-y-[2px]
                  disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? "Creating account..." : "Register"}
              </button>
            </form>
          )}

          {msg && (
            <div className={`mt-4 p-3 text-sm font-bold border-2 border-black rounded-lg
              ${msgType === "success" ? "bg-green-100" : "bg-red-100"}`}
            >
              {msg}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}