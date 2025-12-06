import { error } from "@sveltejs/kit";
import { apiBase, type PostSearchResult } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const search = Array.from(url.searchParams).map(([kind, value]) => ({ kind, value }));
	const res = await fetch(`${apiBase}/search`, { method: "POST", body: JSON.stringify(search) });
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const queryItems: PostSearchResult[] = (await res.json()) ?? [];

	// map from identifier => oldest date for that identifier
	const order = new Map<string, string>();
	for (const { identifier, uploaded } of queryItems) {
		const prev = order.get(identifier);
		if (prev == null || prev > uploaded) {
			// set to the current if there is none
			// or replace if the current uploaded date is older
			order.set(identifier, uploaded);
		}
	}

	const items = queryItems.sort((a, b) => {
		// if the identifiers are unequal, compare
		// the oldest creation dates for each identifier
		// this groups the identifiers together
		if (a.identifier !== b.identifier) {
			const l = order.get(a.identifier)!;
			const r = order.get(b.identifier)!;

			// sort the identifiers in reverse (b cmp a) to
			// put newer identifiers at the top
			return r.localeCompare(l, "en");
		}

		// if they have the same identifier sort their uploaded dates
		// don't sort in reverse (a cmp b) to have ascending creation
		// dates per identifier
		return a.uploaded.localeCompare(b.uploaded, "en");
	});

	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
