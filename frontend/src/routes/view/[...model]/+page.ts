import { apiBase, type Model } from "#lib/api.js";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
	const res = await fetch(`${apiBase}/model?uuid=${encodeURIComponent(params.model)}`);
	if (!res.ok) {
		return error(res.status, await res.json());
	}

	const model: Model = await res.json();
	return { uuid: params.model, model };
};
