import { getLocalModel } from "#lib/local.js";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

// local models are in localStorage, which we don't have access to on the server.
// however, we still want to 404 gracefully when we encounter a model that doesn't exist,
// so we look for it on page load instead of on page mount.
export const ssr = false;

export const load: PageLoad = async ({ params }) => {
	const model = getLocalModel(params.model);
	if (!model) {
		throw error(404, `local model with uuid '${params.model}' not found`);
	}
	return { model };
};
