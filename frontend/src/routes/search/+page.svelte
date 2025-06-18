<script lang="ts">
	import Header from "$lib/Header.svelte";
	import ItemTable from "$lib/ItemTable.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";

	let { data } = $props();
	let { items, params: searchParams } = $derived(data);

	let params = $derived<ChipElement[]>(searchParams.map(([kind, value]) => ({ kind, value })));
</script>

<svelte:head>
	<title>search by metadata</title>
</svelte:head>

<Header {params} />

<ItemTable {items} />
