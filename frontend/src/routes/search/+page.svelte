<script lang="ts">
	import Search from "$lib/Search.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";

	let { data } = $props();
	let { identifiers, params: searchParams } = $derived(data);

	let params = $derived<ChipElement[]>(searchParams.map(([kind, value]) => ({ kind, value })));
</script>

<Search {params} />

<div class="px-4">
	{#each identifiers as { identifier, title }}
		<a class="block text-xl text-blue-500 hover:underline" href="/view/{identifier}">{title}</a>
	{/each}
</div>
