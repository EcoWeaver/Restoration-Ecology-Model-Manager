<script lang="ts">
	import { Search } from "@lucide/svelte";

	let { onclick }: { onclick: () => any } = $props();
</script>

<button
	{onclick}
	title="search"
	class="flex h-10 cursor-pointer place-items-center rounded-r-md border border-l-0 border-neutral-600 px-3 hover:bg-neutral-200 dark:hover:bg-neutral-700"
>
	<Search />
</button>
