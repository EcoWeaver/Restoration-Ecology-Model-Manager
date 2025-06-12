import { apiBase, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/versions${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const items: SearchItem[] = (await res.json()) ?? [];
	const identifier = url.searchParams.get("identifier") || "";
	return { items, identifier };
};
