import { apiBase, sortVersion, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const queryItems: SearchItem[] = (await res.json()) ?? [];

	let order = new Map<string, number>();
	for (let item of queryItems) {
		if (!order.has(item.identifier)) {
			order.set(item.identifier, order.size);
		}
	}

	const items = queryItems.sort((a, b) => {
		let l = order.get(a.identifier)!;
		let r = order.get(b.identifier)!;

		let cmp = l - r;
		if (cmp !== 0) {
			return cmp;
		}

		return sortVersion(a, b);
	});

	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
