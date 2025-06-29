import { apiBase, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const search = Array.from(url.searchParams).map(([kind, value]) => ({ kind, value }));
	const res = await fetch(`${apiBase}/search`, { method: "POST", body: JSON.stringify(search) });
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const items: SearchItem[] = (await res.json()) ?? [];
	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
