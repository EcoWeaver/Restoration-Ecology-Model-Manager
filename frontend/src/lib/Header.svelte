<script lang="ts">
	import Search from "./Search.svelte";
	import type { ChipElement } from "./search/Chip.svelte";

	let { params }: { params?: ChipElement[] } = $props();
</script>

<div class="flex">
	<Search {params} />
	<a
		href="/editor"
		class="m-4 ml-0 flex h-10 place-items-center rounded bg-indigo-300 px-4 hover:bg-indigo-400"
	>
		<span class="font-medium">+ create model</span>
	</a>
</div>
