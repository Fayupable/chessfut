"use client";

import { useEffect, useRef, useState } from "react";
import { searchPlayers } from "@/lib/api";

const PAGE_SIZE = 10;

export function UsernameAutocomplete({
  value,
  onChange,
  placeholder,
}: {
  value: string;
  onChange: (username: string) => void;
  placeholder: string;
}) {
  const [results, setResults] = useState<string[]>([]);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!value.trim()) {
      return;
    }

    const timeout = setTimeout(() => {
      searchPlayers(value.trim(), PAGE_SIZE, 0)
        .then((res) => {
          setResults(res.players);
          setOffset(res.players.length);
          setHasMore(res.players.length === PAGE_SIZE);
          setOpen(true);
        })
        .catch(() => setResults([]));
    }, 700);

    return () => clearTimeout(timeout);
  }, [value]);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const loadMore = () => {
    searchPlayers(value.trim(), PAGE_SIZE, offset).then((res) => {
      setResults((prev) => [...prev, ...res.players]);
      setOffset((prev) => prev + res.players.length);
      setHasMore(res.players.length === PAGE_SIZE);
    });
  };

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const el = e.currentTarget;
    if (hasMore && el.scrollHeight - el.scrollTop - el.clientHeight < 40) {
      loadMore();
    }
  };

  const showDropdown = open && value.trim().length > 0 && results.length > 0;

  return (
    <div ref={containerRef} className="relative">
      <input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onFocus={() => value.trim() && setOpen(true)}
        placeholder={placeholder}
        className="w-full rounded-md border border-neutral-700 bg-neutral-800 px-4 py-2 outline-none"
      />
      {showDropdown && (
        <div
          onScroll={handleScroll}
          className="absolute z-20 mt-1 max-h-56 w-full overflow-y-auto rounded-md border border-white/10 bg-neutral-800 shadow-lg"
        >
          {results.map((username) => (
            <button
              key={username}
              type="button"
              onClick={() => {
                onChange(username);
                setOpen(false);
              }}
              className="block w-full px-3 py-2 text-left text-sm text-white hover:bg-neutral-700"
            >
              {username}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}