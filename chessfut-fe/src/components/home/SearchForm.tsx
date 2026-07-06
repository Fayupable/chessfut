"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { UsernameAutocomplete } from "@/components/search/UsernameAutocomplete";

export function SearchForm() {
  const [username, setUsername] = useState("");
  const router = useRouter();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      router.push(`/${username.trim()}`);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <div className="flex-1">
        <UsernameAutocomplete value={username} onChange={setUsername} placeholder="chess.com username" />
      </div>
      <button
        type="submit"
        className="rounded-md bg-emerald-500 px-4 py-2 font-semibold text-black hover:bg-emerald-400"
      >
        Scout
      </button>
    </form>
  );
}