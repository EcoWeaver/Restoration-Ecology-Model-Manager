import { apiBase } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, url }) => {
	const res = await fetch(`${apiBase}/search${url.search}`);
	const identifiers: string[] = await res.json();

	return { identifiers };
};
