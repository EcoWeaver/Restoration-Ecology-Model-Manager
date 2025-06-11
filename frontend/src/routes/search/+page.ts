import { apiBase, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const items: SearchItem[] = await res.json();
	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
