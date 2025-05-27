import { apiBase } from "$lib/api";
import type { PageLoad } from "./$types";

type Item = {
	identifier: string;
	title: string;
};

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	const identifiers: Item[] = await res.json();

	let params = Array.from(url.searchParams.entries());
	return { identifiers, params };
};
