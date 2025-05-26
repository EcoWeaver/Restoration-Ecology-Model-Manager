<script lang="ts">
	let { data } = $props();
	let { identifiers } = $derived(data);

	import Search from "$lib/Search.svelte";
</script>

<Search />

<div class="px-4">
	{#each identifiers as { identifier, title }}
		<a class="block text-xl text-blue-500 hover:underline" href="/view/{identifier}">{title}</a>
	{/each}
</div>
