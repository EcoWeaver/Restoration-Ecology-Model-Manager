<script lang="ts">
	import Header from "$lib/Header.svelte";
	import ItemTable from "$lib/ItemTable.svelte";

	let { data } = $props();
	let { items, identifier } = $derived(data);
</script>

<Header {identifier} />

<ItemTable {items} />
