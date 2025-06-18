<script lang="ts">
	import Header from "$lib/Header.svelte";
	import ItemTable from "$lib/ItemTable.svelte";

	let { data } = $props();
	let { items, identifier } = $derived(data);
</script>

<svelte:head>
	<title>search by identifier</title>
</svelte:head>

<Header {identifier} />

<ItemTable {items} />
