<script lang="ts">
	import { CodeXml, House, Plus } from "@lucide/svelte";
	import Search from "./Search.svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import type { Snippet } from "svelte";
	import DarkLight from "./DarkLight.svelte";

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

	<div class="flex min-w-0 grow">
		<Search {params} {identifier} selectedSearch={identifier != null ? "identifier" : "metadata"} />
	</div>

	{#if createModel}
		<a
			href="/edit"
			title="create new model"
			class="flex h-10 items-center rounded-md bg-indigo-300 px-2 font-medium hover:bg-indigo-400 dark:bg-indigo-500"
		>
			<Plus />
			<span class="px-1">model</span>
		</a>
	{/if}

	{@render children?.()}

	<DarkLight />

	<a
		title="source"
		href="https://git.uni-jena.de/fusion/teaching/project/2025sose/swep/restoration-ecology-model-manager/restoration-ecology-model-manager-app"
		target="_blank"
		rel="noopener noreferrer"
	>
		<CodeXml />
	</a>
</div>
