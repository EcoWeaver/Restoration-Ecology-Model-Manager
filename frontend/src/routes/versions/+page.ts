import { apiBase, sortVersion, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/versions${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const queryItems: SearchItem[] = (await res.json()) ?? [];
	const items = queryItems.sort(sortVersion);

	const identifier = url.searchParams.get("identifier") || "";
	return { items, identifier };
};
