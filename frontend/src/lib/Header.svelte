<script lang="ts">
	import { CodeXml, House } from "@lucide/svelte";
	import Search from "./Search.svelte";
	import type { ChipElement } from "./search/Chip.svelte";

	let { params, identifier }: { params?: ChipElement[]; identifier?: string } = $props();
</script>

<div class="flex">
	<a class="m-4 mr-0 pt-1 pb-1" title="home" href="/">
		<House class="size-8" />
	</a>

	<Search {params} {identifier} selectedSearch={identifier != null ? "ident" : "meta"} />

	<a
		href="/edit"
		class="m-4 ml-0 flex h-10 place-items-center rounded-md bg-indigo-300 px-4 hover:bg-indigo-400 dark:bg-indigo-500"
	>
		<span class="font-medium">+ create model</span>
	</a>

	<a
		class="m-4 ml-0 pt-1 pb-1"
		title="source"
		href="https://git.uni-jena.de/fusion/teaching/project/2025sose/swep/restoration-ecology-model-manager/restoration-ecology-model-manager-app"
		target="_blank"
		rel="noopener noreferrer"
	>
		<CodeXml class="size-8" />
	</a>
</div>
