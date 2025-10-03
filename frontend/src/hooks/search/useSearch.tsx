import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";

export function useSearch() {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.code === "Space") {
        e.preventDefault();
        setOpen(true);
      }
      if (e.key === "Escape") {
        setOpen(false);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const SearchModal = () =>
    open ? (
      <div className="fixed inset-0 z-50 bg-black/70 flex items-center justify-center">
        <div
          className="bg-white text-black border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] p-8 text-center"
          style={{
           
            fontFamily: "'Space Grotesk', sans-serif",
          }}
        >
          <img
            src="/undraw/not-done.svg"
            alt="MacBook"
            className="w-56 mx-auto mb-8 rounded-xl border-4 border-black"
          />
          <h2 className="text-3xl font-extrabold mb-6 tracking-tight">
            Search Coming Soon
          </h2>
          <Button
            variant="noShadow"
            size="lg"
            className="bg-red-600 text-white border-2 border-black rounded-lg px-4 py-2 font-bold shadow-[4px_4px_0_0_#000] hover:bg-red-700 transition-all"
            onClick={() => setOpen(false)}
          >
            Close
          </Button>
        </div>
      </div>
    ) : null;

  return { open, setOpen, SearchModal };
}