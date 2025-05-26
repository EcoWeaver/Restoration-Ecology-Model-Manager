import { apiBase, type ModelRequest } from "$lib/api";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ fetch, params }) => {
	const res = await fetch(`${apiBase}/model?identifier=${params.model}`);
	const model: ModelRequest = await res.json();

	return { model };
};
