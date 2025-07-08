<script lang="ts">
	import { CodeXml, House, Plus } from "@lucide/svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import type { Snippet } from "svelte";
	import DarkLight from "./DarkLight.svelte";
	import MetadataSearch from "./MetadataSearch.svelte";

	let {
		params,
		children,
		createModel = true,
	}: {
		params?: ChipElement[];
		createModel?: boolean;
		children?: Snippet;
	} = $props();
</script>

<div class="flex items-center gap-4 p-4" role="navigation">
	<a
		title="home"
		class="rounded-4xl outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		href="/"
	>
		<House />
	</a>

	<div class="flex min-w-0 grow">
		<MetadataSearch {params} />
	</div>

	{#if createModel}
		<a
			href="/edit"
			title="create new model"
			class="flex h-10 items-center rounded-md bg-indigo-300 px-2 font-medium outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
		>
			<Plus />
			<span class="px-1">model</span>
		</a>
	{/if}

	{@render children?.()}

	<DarkLight />

	<a
		class="rounded-full outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		title="source code"
		href="https://git.uni-jena.de/fusion/teaching/project/2025sose/swep/restoration-ecology-model-manager/restoration-ecology-model-manager-app"
		target="_blank"
		rel="noopener noreferrer"
	>
		<CodeXml />
	</a>
</div>
