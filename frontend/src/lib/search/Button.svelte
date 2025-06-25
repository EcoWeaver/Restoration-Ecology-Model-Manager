<script lang="ts">
	import type { MouseEventHandler } from "svelte/elements";
	import { Search } from "@lucide/svelte";

	let { onclick }: { onclick: MouseEventHandler<HTMLButtonElement> } = $props();
</script>

<button
	{onclick}
	title="search"
	class="flex h-10 cursor-pointer place-items-center rounded-r-md border border-l-0 border-neutral-600 px-3 outline-indigo-400 focus-visible:outline-2"
>
	<Search />
</button>
