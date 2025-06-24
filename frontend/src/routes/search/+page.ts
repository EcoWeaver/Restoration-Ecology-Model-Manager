import { apiBase, sortVersion, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const queryItems: SearchItem[] = (await res.json()) ?? [];

	const order = new Map<string, number>();
	for (const item of queryItems) {
		if (!order.has(item.identifier)) {
			order.set(item.identifier, order.size);
		}
	}

	const items = queryItems.sort((a, b) => {
		const l = order.get(a.identifier)!;
		const r = order.get(b.identifier)!;

		const cmp = l - r;
		if (cmp !== 0) {
			return cmp;
		}

		return sortVersion(a, b);
	});

	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
