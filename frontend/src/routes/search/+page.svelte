<script lang="ts">
	import Header from "$lib/Header.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";
	import SearchResults from "$lib/SearchResults.svelte";

	let { data } = $props();

	let { items, params: searchParams } = $derived(data);

	let params = $derived<ChipElement[]>(searchParams.map(([kind, value]) => ({ kind, value })));
</script>

<svelte:head>
	<title>search uploaded models</title>
</svelte:head>

<Header {params} />

<SearchResults {items} route="view" />
