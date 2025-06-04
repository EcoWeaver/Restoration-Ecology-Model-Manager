import { apiBase } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

type Item = {
	identifier: string;
	version: string,
	title: string;
};

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const items: Item[] = await res.json();

	const params = Array.from(url.searchParams.entries());
	return { items, params };
};
