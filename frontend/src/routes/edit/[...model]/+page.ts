import { apiBase, type GetModelResult } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
	const apiUrl = new URL("/model", apiBase);
	apiUrl.searchParams.set("uuid", params.model);

	const res = await fetch(apiUrl);
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const model: GetModelResult = await res.json();
	return { model };
};
