import React, { useState } from "react";
import { NeoButton } from "@/components/ui/neo-button";
import http from "@/http";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog";

export default function CreateNotebook({ onCreated }: { onCreated?: () => void }) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [tags, setTags] = useState<string>("");
  const [color, setColor] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleCreate = async () => {
    setLoading(true);
    setError(null);
    try {
      await http.post("/api/v1/notebooks", {
        name,
        description,
        tags: tags.split(",").map(t => t.trim()).filter(Boolean),
        color,
      });
      setOpen(false);
      setName("");
      setDescription("");
      setTags("");
      setColor("");
      if (onCreated) onCreated();
    } catch (err: any) {
      setError(err?.response?.data?.error || "Failed to create notebook");
    }
    setLoading(false);
  };

  return (
    <>
      <NeoButton color="green" title="Create Notebook" onClick={() => setOpen(true)} />
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Notebook</DialogTitle>
          </DialogHeader>
          <div className="mb-3">
            <label className="block font-semibold mb-1">Name</label>
            <input
              className="w-full border-2 border-black rounded px-3 py-2"
              value={name}
              onChange={e => setName(e.target.value)}
              placeholder="Notebook name"
            />
          </div>
          <div className="mb-3">
            <label className="block font-semibold mb-1">Description</label>
            <input
              className="w-full border-2 border-black rounded px-3 py-2"
              value={description}
              onChange={e => setDescription(e.target.value)}
              placeholder="Notebook description"
            />
          </div>
          <div className="mb-3">
            <label className="block font-semibold mb-1">Tags (comma separated)</label>
            <input
              className="w-full border-2 border-black rounded px-3 py-2"
              value={tags}
              onChange={e => setTags(e.target.value)}
              placeholder="math, science, ai"
            />
          </div>
          <div className="mb-3">
            <label className="block font-semibold mb-1">Color</label>
            <input
              className="w-full border-2 border-black rounded px-3 py-2"
              value={color}
              onChange={e => setColor(e.target.value)}
              placeholder="blue, red, etc"
            />
          </div>
          {error && <div className="text-red-600 mb-2">{error}</div>}
          <DialogFooter>
            <NeoButton
              color="green"
              title={loading ? "Creating..." : "Create"}
              onClick={handleCreate}
              disabled={loading || !name}
              className="w-full"
            />
            <DialogClose asChild>
              <NeoButton color="red" title="Cancel" className="w-full mt-2" />
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}