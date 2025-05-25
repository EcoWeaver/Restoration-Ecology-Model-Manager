import { apiBase } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
	const res = await fetch(`${apiBase}/model?identifier=${params.model}`)
	const model = await res.json();

	return { model };
};
