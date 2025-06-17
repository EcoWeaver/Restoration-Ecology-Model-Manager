import { apiBase, type ModelRequest } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params, url }) => {
	let apiUrl = new URL("/model", apiBase);
	apiUrl.searchParams.set("identifier", params.model);
	apiUrl.searchParams.set("version", url.searchParams.get("version") ?? "");

	const res = await fetch(apiUrl);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const model: ModelRequest = await res.json();
	return { model };
};
