<script lang="ts">
	import { Link, PencilLine } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import Header from "$lib/Header.svelte";
	import Alert from "$lib/Alert.svelte";
	import { apiToFlow } from "$lib/model";

	let { data } = $props();
	let { model, metadataMap } = $derived(data);
	let { nodes, edges } = $derived(apiToFlow(model.graph, false));

	let alert: ReturnType<typeof Alert>;
	async function share() {
		await navigator.clipboard.writeText(window.location.href);
		alert.info("link copied to clipboard");
	}
</script>

<svelte:head>
	<title>view {model.identifier}@v{model.version}</title>
</svelte:head>

<div class="absolute z-10 w-full bg-white dark:bg-neutral-800">
	<Header>
		<a
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="edit model"
			href={`/edit/${model.identifier}?version=${model.version}`}
		>
			<PencilLine strokeWidth={1.7} />
			<span class="pl-2">edit</span>
		</a>

		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="copy model link to clipboard"
			onclick={share}
		>
			<Link />
			<span class="pl-2">share</span>
		</button>
	</Header>
</div>

<Alert bind:this={alert} />

<div class="h-dvh w-dvw">
	<Graph {nodes} {edges} interactive={false}>
		{#snippet panel()}
			<PanelComponent>
				<h2 class="text-xl">Metadata</h2>

				<div class="table">
					<div class="table-row">
						<div class="table-cell pr-2 pb-1 font-bold">kind</div>
						<div class="table-cell pb-1 pl-2 font-bold">value</div>
					</div>
					<div class="table-row">
						<div class="table-cell pr-2">identifier</div>
						<a
							class="max-2-128 table-cell pl-2 text-blue-700 outline-indigo-400 hover:underline focus-visible:outline-2 dark:text-blue-400"
							href={`https://doi.org/${model.identifier}`}
							target="_blank"
							rel="noopener noreferrer"
						>
							{model.identifier}
						</a>
					</div>
					{#each model.metadata as { kind, value }, idx (idx)}
						{@const type = metadataMap.get(kind)}

						<div class="table-row">
							<div class="table-cell pr-2">{kind}</div>

							{#if type === "orcid"}
								<a
									class="table-cell max-w-128 pl-2 text-blue-700 outline-indigo-400 hover:underline focus-visible:outline-2 dark:text-blue-400"
									href={`https://orcid.org/${value}`}
									target="_blank"
									rel="noopener noreferrer"
								>
									{value}
								</a>
							{:else if type === "url"}
								<a
									class="table-cell max-w-128 pl-2 text-blue-700 hover:underline dark:text-blue-400"
									href={value}
									target="_blank"
									rel="noopener noreferrer"
								>
									{value}
								</a>
							{:else}
								<div class="table-cell max-w-128 pl-2">{value}</div>
							{/if}
						</div>
					{/each}
				</div>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>
