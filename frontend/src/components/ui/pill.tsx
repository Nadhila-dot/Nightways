import React from "react";
import { useEffect, useState } from "react";

interface PillProps {
  text: string;
  action?: string;
  className?: string;
  onClick?: () => void;
}

const Pill: React.FC<PillProps> = ({ text, action, className, onClick }) => {
  const [platform, setPlatform] = useState<"mac" | "windows" | "mobile" | "other">("other");

  useEffect(() => {
    // Detect platform
    const userAgent = navigator.userAgent.toLowerCase();
    if (/iphone|ipad|ipod|macintosh/.test(userAgent)) {
      setPlatform("mac");
    } else if (/windows/.test(userAgent)) {
      setPlatform("windows");
    } else if (/android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/.test(userAgent)) {
      setPlatform("mobile");
    } else {
      setPlatform("other");
    }
  }, []);

  const renderIcon = () => {
    switch (platform) {
      case "mac":
        return (
          <span className="inline-block bg-white text-black font-mono border-2 border-black rounded px-1 shadow-[2px_2px_0_0_#000] text-xs">
            ⇥
          </span>
        );
      case "windows":
        return (
          <span className="inline-block bg-white text-black font-mono border-2 border-black rounded px-1 shadow-[2px_2px_0_0_#000] text-xs">
            Tab
          </span>
        );
      case "mobile":
        return (
          <span className="inline-block bg-black text-white rounded-full w-5 h-5 text-xs flex items-center justify-center mr-1">
            ✨
          </span>
        );
      default:
        return (
          <span className="inline-block bg-white text-black font-mono border-2 border-black rounded px-1 shadow-[2px_2px_0_0_#000] text-xs">
            Tab
          </span>
        );
    }
  };

  return (
    <div
      className={`
        inline-flex items-center px-3 py-2  
        bg-yellow-300 text-black font-medium
        border-2 border-black rounded-md  
        shadow-[3px_3px_0_0_#000] 
        hover:translate-x-[1px] hover:translate-y-[1px] hover:shadow-[2px_2px_0_0_#000]
        active:translate-x-[3px] active:translate-y-[3px] active:shadow-none
        transition-all cursor-pointer animate - later something
        ${className || ""}
      `}
      onClick={onClick}
    >
      <div className="flex items-center gap-1.5">
        {renderIcon()}
        <span className="text-sm font-bold">{text}</span>
        {action && <span className="text-xs opacity-80">to {action}</span>}
      </div>
    </div>
  );
};

export { Pill };