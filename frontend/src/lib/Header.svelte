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

<div class="flex items-center gap-4 p-4" role="navigation">
	<a
		title="home"
		class="rounded-4xl outline-offset-4 outline-indigo-400 focus-visible:outline-2"
		href={resolve("/")}
	>
		<House />
	</a>

	<div class="flex min-w-0 grow">
		<SearchBar {params} />
	</div>

	{#if yourModels}
		<a
			href={resolve("/edit")}
			title="view your local models"
			class="flex h-10 items-center rounded-md bg-indigo-300 px-3 font-medium outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
		>
			<List />
			<span class="pl-2 select-none">your models</span>
		</a>
	{/if}

	{#if newModel}
		<button
			onclick={() => goto(resolve(`/edit/${newLocalModel()}`))}
			title="create new model"
			class="flex h-10 items-center rounded-md bg-indigo-300 px-3 font-medium outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
		>
			<Plus />
			<span class="pl-2 select-none">model</span>
		</button>
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
