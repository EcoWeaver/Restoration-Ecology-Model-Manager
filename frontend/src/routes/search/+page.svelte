<script lang="ts">
	let { data } = $props();
	let { identifiers } = $derived(data);

	import Search from "$lib/Search.svelte";
</script>

<Search />

<div class="p-2">
	{#each identifiers as item}
		<a class="block text-blue-500 hover:underline" href="/view/{item}">{item}</a>
	{/each}
</div>
