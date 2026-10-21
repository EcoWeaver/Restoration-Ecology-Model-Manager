import { error } from "@sveltejs/kit";
import { apiBase, type PostSearchResult } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const metadata = Array.from(url.searchParams).map(([kind, value]) => ({ kind, value }));
	const res = await fetch(`${apiBase}/search`, {
		method: "POST",
		body: JSON.stringify(metadata),
	});
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const items: PostSearchResult[] = (await res.json()) ?? [];
	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
