<script lang="ts">
	import Header from "#lib/Header.svelte";
	import { deleteLocalModel, getLocalModels } from "#lib/local.js";
	import SearchResults from "#lib/SearchResults.svelte";

	function placeholder<T extends Record<string, string>>(obj: T): T {
		return Object.fromEntries(
			Object.entries(obj).map(([key, value]) => [key, value === "" ? "n/a" : value]),
		) as T;
	}

	let localModels = $state(getLocalModels());
	const items = $derived(
		localModels
			.entries()
			.map(([uuid, res]) => ({ uuid, ...res }))
			.map((res) => placeholder(res))
			.toArray(),
	);
</script>

<svelte:head>
	<title>view your local models</title>
</svelte:head>

<Header newModel={true} yourModels={false} />

<SearchResults
	{items}
	route="edit"
	del={(uuid) => {
		deleteLocalModel(uuid);
		localModels = getLocalModels();
	}}
/>
