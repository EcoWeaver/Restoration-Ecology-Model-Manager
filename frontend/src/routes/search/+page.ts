import { apiBase, type SearchItem } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const search = Array.from(url.searchParams).map(([kind, value]) => ({ kind, value }));
	const res = await fetch(`${apiBase}/search`, { method: "POST", body: JSON.stringify(search) });
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const queryItems: SearchItem[] = (await res.json()) ?? [];

	const order = new Map<string, string>();
	for (const { identifier, created } of queryItems) {
		const prev = order.get(identifier);
		if (prev == null || prev > created) {
			order.set(identifier, created);
		}
	}

	const items = queryItems.sort((a, b) => {
		const l = order.get(a.identifier)!;
		const r = order.get(b.identifier)!;

		if (r < l) {
			return -1;
		} else if (l < r) {
			return 1;
		}

		if (a.identifier < b.identifier) {
			return -1;
		} else if (b.identifier < a.identifier) {
			return 1;
		} else {
			return 0;
		}
	});

	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
