import { apiBase, type ModelRequest } from "$lib/api";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
	const res = await fetch(`${apiBase}/model?identifier=${params.model}`);
	if (!res.ok) {
		return error(res.status, await res.text());
	}

	const model: ModelRequest = await res.json();

	return { model };
};
