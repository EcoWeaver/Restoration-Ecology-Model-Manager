import { apiBase, sortVersion, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const queryItems: SearchItem[] = (await res.json()) ?? [];
	const items = queryItems.sort((a, b) => {
		if (a.identifier < b.identifier) {
			return -1;
		} else if (a.identifier > b.identifier) {
			return 1;
		}

		return sortVersion(a, b);
	});

	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
