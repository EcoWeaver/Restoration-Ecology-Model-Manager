import { apiBase } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

type Item = {
	identifier: string;
	title: string;
};

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const identifiers: Item[] = await res.json();

	let params = Array.from(url.searchParams.entries());
	return { identifiers, params };
};
