<script lang="ts">
	import { CodeXml, House, Plus } from "@lucide/svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import type { Snippet } from "svelte";
	import DarkLight from "./DarkLight.svelte";
	import MetadataSearch from "./MetadataSearch.svelte";
	import { type MetadataElements } from "./api";
	import { resolve } from "$app/paths";

	let {
		metadataElements,
		params,
		children,
		createModel = true,
	}: {
		metadataElements: MetadataElements;
		params?: ChipElement[];
		createModel?: boolean;
		children?: Snippet;
	} = $props();
</script>

<div class="flex items-center gap-4 p-4" role="navigation">
	<a
		title="home"
		class="rounded-4xl outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		href={resolve("/")}
	>
		<House />
	</a>

	<div class="flex min-w-0 grow">
		<MetadataSearch {metadataElements} {params} />
	</div>

	{#if createModel}
		<a
			href={resolve("/edit")}
			title="create new model"
			class="flex h-10 items-center rounded-md bg-indigo-300 px-2 font-medium outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
		>
			<Plus />
			<span class="px-1 select-none">model</span>
		</a>
	{/if}

	{@render children?.()}

	<DarkLight />

	<a
		class="rounded-full outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		title="source code"
		href="https://github.com/EcoWeaver/Restoration-Ecology-Model-Manager"
		target="_blank"
		rel="noopener noreferrer"
	>
		<CodeXml />
	</a>
</div>
