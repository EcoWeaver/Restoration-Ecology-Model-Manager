<script lang="ts">
	import { CodeXml, House, List, Plus } from "@lucide/svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import type { Snippet } from "svelte";
	import DarkLight from "./DarkLight.svelte";
	import SearchBar from "./SearchBar.svelte";
	import { resolve } from "$app/paths";
	import { goto } from "$app/navigation";
	import { newLocalModel } from "./local";

	let {
		params,
		newModel = false,
		yourModels = true,
		children,
	}: {
		params?: ChipElement[];
		newModel?: boolean;
		yourModels?: boolean;
		children?: Snippet;
	} = $props();
</script>

<div class="flex items-center gap-4 p-4 whitespace-nowrap" role="navigation">
	<a
		class="rounded-full outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		title="home"
		href={resolve("")}
	>
		<House />
	</a>

	<div class="flex min-w-0 grow">
		<SearchBar {params} />
	</div>

	{#if yourModels}
		<a
			class="flex h-10 items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="view your local models"
			href={resolve("edit")}
		>
			<List />
			<div class="pl-2 select-none">your models</div>
		</a>
	{/if}

	{#if newModel}
		<button
			class="flex h-10 items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="create new model"
			onclick={() => goto(resolve(`edit/${newLocalModel()}`))}
		>
			<Plus />
			<div class="pl-2">model</div>
		</button>
	{/if}

	{@render children?.()}

	<DarkLight />

	<a
		class="rounded-full outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		title="source code"
		href="https://github.com/EcoWeaver/Restoration-Ecology-Model-Manager"
		target="_blank"
		rel="external noreferrer"
	>
		<CodeXml />
	</a>
</div>
