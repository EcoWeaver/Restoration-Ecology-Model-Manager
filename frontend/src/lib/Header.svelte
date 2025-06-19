<script lang="ts">
	import { CodeXml, House, Plus } from "@lucide/svelte";
	import Search from "./Search.svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import type { Snippet } from "svelte";

	let {
		params,
		identifier,
		children,
		createModel = true,
	}: {
		params?: ChipElement[];
		identifier?: string;
		createModel?: boolean;
		children?: Snippet;
	} = $props();
</script>

<div class="flex items-center gap-4 p-4" role="navigation">
	<a title="home" href="/"><House /></a>

	<div class="flex grow">
		<Search {params} {identifier} selectedSearch={identifier != null ? "ident" : "meta"} />
	</div>

	{#if createModel}
		<a
			href="/edit"
			class="flex h-10 items-center rounded-md bg-indigo-300 px-2 font-medium hover:bg-indigo-400 dark:bg-indigo-500"
		>
			<Plus />
			<span class="px-1">new model</span>
		</a>
	{/if}

	{@render children?.()}

	<a
		title="source"
		href="https://git.uni-jena.de/fusion/teaching/project/2025sose/swep/restoration-ecology-model-manager/restoration-ecology-model-manager-app"
		target="_blank"
		rel="noopener noreferrer"
	>
		<CodeXml />
	</a>
</div>
