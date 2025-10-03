import { useEffect, useState } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";
import { Clock } from "lucide-react";
import { useAuth } from "@/hooks/auth/checkAuth";
import http from "@/http";
import { getUserInfo } from "@/scripts/getUserinfo";
import { useSession } from "@/hooks/auth/getSession";
import { getSystemInfo } from "@/scripts/getSystem";
import { NeoButton } from "../ui/neo-button";
import { toast } from "sonner";


type UserInfo = {
  username: string;
  email: string;
  rank: string;
  avatar?: string | null;
};



export default function DashboardCard() {
  const isAuthenticated = useAuth();
  const session = useSession();
  const [user, setUser] = useState<UserInfo | null>(null);
  const [systemInfo, setSystemInfo] = useState<any>(null);
 


   useEffect(() => {
    getSystemInfo().then((info) => setSystemInfo(info));
  }, []);

   useEffect(() => {
    if (session) {
        getUserInfo(session).then((data) => setUser(data));
    }
    }, [session]);
  const [currentTime, setCurrentTime] = useState<string>("");

  useEffect(() => {
    const updateTime = () => {
      const now = new Date();
      setCurrentTime(
        now.toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        })
      );
    };
    updateTime();
    const interval = setInterval(updateTime, 1000);
    return () => clearInterval(interval);
  }, []);

  if (!user) return null;

  

  if (!user || !systemInfo) return null;

  return (
  <div className="w-full max-w-full mx-auto mt-8 bg-white border-4 border-black shadow-[8px_8px_0px_0px_#000000] rounded-lg">
    <div className="flex items-center gap-6 px-12 py-8">
      <div className="flex items-center gap-4">
        <div className="flex flex-col items-center">
          <div className="flex items-center gap-2 bg-yellow-300 border-3 border-black px-3 py-1 shadow-[3px_3px_0px_0px_#000000] rounded-md justify-start self-start">
            <Clock className="h-5 w-5 text-black" />
            <span className="text-3xl font-black text-black">{currentTime}</span>
          </div>
          <p className="mt-2 font-bold text-black text-left self-start">Build: {systemInfo.data.build}</p>
        </div>
        
        
      </div>
      <div className="flex-1 text-right">
        <div className="text-3xl font-black text-black mb-1 uppercase tracking-tight">{user.username}</div>
        <div className="text-lg font-bold text-gray-800 mb-4">{user.email}</div>
        <div className="flex justify-end">
          {/*<span className="px-6 py-2 bg-lime-300 border-4 border-black font-black text-lg text-black uppercase shadow-[4px_4px_0px_0px_#000000] transform -rotate-1 rounded-md">
            {user.rank}
          </span>*/}
          <NeoButton
          title={user.rank.toUpperCase()}
          // fancy ass code man
          color={user.rank === "admin" ? "red" : user.rank === "moderator" ? "yellow" : "green"}
          onClick={() => {
            toast.success("Your rank is " + user.rank + " :D");
          }}
          />
        </div>
      </div>
    </div>
  </div>
);
}