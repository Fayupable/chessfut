"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { UsernameAutocomplete } from "@/components/search/UsernameAutocomplete";

export function CompareSearchForm() {
    const [usernameA, setUsernameA] = useState("");
    const [usernameB, setUsernameB] = useState("");
    const router = useRouter();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (usernameA.trim() && usernameB.trim()) {
            router.push(`/compare?a=${usernameA.trim()}&b=${usernameB.trim()}`);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="flex flex-col items-center gap-4">
            <div className="flex items-center gap-3">
                <UsernameAutocomplete value={usernameA} onChange={setUsernameA} placeholder="First username" />
                <span className="text-white/40">vs</span>
                <UsernameAutocomplete value={usernameB} onChange={setUsernameB} placeholder="Second username" />
            </div>
            <button
                type="submit"
                className="rounded-md bg-emerald-500 px-6 py-2 font-semibold text-black hover:bg-emerald-400"
            >
                Compare
            </button>
        </form>
    );
}