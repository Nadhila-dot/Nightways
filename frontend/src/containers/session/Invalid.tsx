import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

export default function InvalidSessionModal() {
  const navigate = useNavigate();

  // Prevent closing with Escape or click outside
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        e.preventDefault();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  return (
    <Dialog open={true} modal={true}>
      <DialogContent
        className="sm:max-w-[400px] bg-white border-4 border-black rounded-xl shadow-[8px_8px_0_0_#000] pointer-events-auto"
       // hideClose
      >
        <DialogHeader>
          <DialogTitle className="text-3xl font-extrabold text-black tracking-tight">
            Session Invalid
          </DialogTitle>
          <DialogDescription className="text-base text-black">
            Your session has expired or is invalid.<br />
            Please log in again to continue.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="mt-6 flex justify-center">
          <Button
            className="px-6 py-3 font-extrabold text-black hover:text-white hover:bg-black hover:shadow-2xl  transition-all"
            onClick={() => {
                navigate("/auth/login");
                window.location.reload();
            }}
          >
            Go to Login
          </Button>
        </DialogFooter>
      </DialogContent>
      {/* Prevent interaction with background */}
      <div className="fixed inset-0 bg-black bg-opacity-60 z-40 pointer-events-auto" />
    </Dialog>
  );
}