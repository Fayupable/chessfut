"use client";

import { useState } from "react";

export function CardActions({ username }: { username: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    const url = `${window.location.origin}/${username}`;
    await navigator.clipboard.writeText(url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="flex w-[300px] flex-col gap-2">
      <button
        onClick={handleCopy}
        className="rounded-md bg-emerald-500 px-4 py-2 font-semibold text-black hover:bg-emerald-400"
      >
        {copied ? "Link copied!" : "Copy link"}
      </button>
    </div>
  );
}